# kratos概览

> [kratos](github.com/go-kratos/kratos) is a caching and cache-filling library, intended as a replacement for memcached in many cases.

## 介绍

实现一个简单的http server，使用mysql和memcache

## 创建数据库表

```sql
CREATE TABLE `brand` (
	`id` BIGINT NOT NULL AUTO_INCREMENT,
	`name` VARCHAR(64) NULL DEFAULT NULL,
	`first_letter` VARCHAR(8) NULL DEFAULT NULL COMMENT '首字母',
	`sort` INT NULL DEFAULT NULL,
	`factory_status` INT NULL DEFAULT NULL COMMENT '是否为品牌制造商：0->不是；1->是',
	`show_status` INT NULL DEFAULT NULL,
	`product_count` INT NULL DEFAULT NULL COMMENT '产品数量',
	`product_comment_count` INT NULL DEFAULT NULL COMMENT '产品评论数量',
	`logo` VARCHAR(255) NULL DEFAULT NULL COMMENT '品牌logo',
	`big_pic` VARCHAR(255) NULL DEFAULT NULL COMMENT '专区大图',
	`brand_story` TEXT NULL COMMENT '品牌故事',
	PRIMARY KEY (`id`)
)
COMMENT='品牌表'
COLLATE='utf8_general_ci'
ENGINE=InnoDB
AUTO_INCREMENT=59
;
```

## 创建项目并定义proto

### 创建项目

`kratos new brand`

### 定义proto

```proto3
syntax = "proto3";

import "github.com/gogo/protobuf/gogoproto/gogo.proto";
import "google/protobuf/empty.proto";
import "google/api/annotations.proto";

package demo.service.v1;

option go_package = "api";
option (gogoproto.goproto_getters_all) = false;

service Demo {
  rpc Ping(.google.protobuf.Empty) returns (.google.protobuf.Empty);
  rpc Brands(.google.protobuf.Empty) returns (BrandsReply) {
    option (google.api.http) = {
      get: "/api/v1/brands"
    };
  };
}

message Brand {
  int64 id = 1;
  string name = 2;
  string first_letter = 3;
  int64 sort = 4;
  int64 factory_status = 5;
  int64 show_status = 6;
  int64 product_count = 7;
  int64 product_comment_count = 8;
  string logo = 9;
  string big_pic = 10;
  string brand_story = 11;
}

message BrandsReply {
  repeated Brand brands = 1;
}
```

## 在service包中实现grpc server

```go
// Ping ping the resource.
func (s *Service) Ping(ctx context.Context, e *empty.Empty) (*empty.Empty, error) {
	return &empty.Empty{}, s.dao.Ping(ctx)
}

// Close close the resource.
func (s *Service) Close() {
	s.dao.Close()
}

// Brands return all brands
func (s *Service) Brands(ctx context.Context, e *empty.Empty) (reply *pb.BrandsReply, err error) {
	mid, err := s.dao.BrandList(ctx)
	if err != nil {
		return
	}
	brandsMap, err := s.dao.BrandInfoList(ctx, mid)
	if err != nil {
		return nil, err
	}
	var brands []*pb.Brand
	for _, v := range brandsMap {
		brands = append(brands, v)
	}
	reply = &pb.BrandsReply{
		Brands: brands,
	}
	return
}
```

```go
type Dao interface {
	Close()
	Ping(ctx context.Context) (err error)
	// bts: -nullcache=&pb.Brand{Id:-1} -check_null_code=$!=nil&&$.Id==-1
	BrandInfoList(c context.Context, ids []int64) (map[int64]*pb.Brand, error) // 缓存此方法，批量缓存
	BrandList(c context.Context) ([]int64, error) // 此方法不缓存
}
```

## 定义dao和缓存回源并实现

```go
type Dao interface {
	Close()
	Ping(ctx context.Context) (err error)
	// bts: -nullcache=&pb.Brand{Id:-1} -check_null_code=$!=nil&&$.Id==-1
	BrandInfoList(c context.Context, ids []int64) (map[int64]*pb.Brand, error) // 缓存此方法，批量缓存
	BrandList(c context.Context) ([]int64, error) // 此方法不缓存
}

// Close close the resource.
func (d *dao) Close() {
	d.cache.Close()
	d.mc.Close()
}

// Ping ping the resource.
func (d *dao) Ping(ctx context.Context) (err error) {
	return d.pingMC(ctx)
}

const (
	_brandInfoList = "SELECT `id`, IFNULL(`name`, ''), IFNULL(`first_letter`, ''), IFNULL(`sort`, 0), IFNULL(`factory_status`, 0), IFNULL(`show_status`, 0), IFNULL(`product_count`, 0), IFNULL(`product_comment_count`, 0), IFNULL(`logo`, ''), IFNULL(`big_pic`, ''), IFNULL(`brand_story`, '') FROM `pms_brand` WHERE `id` in (%s)"
	_brandList     = "SELECT `id` from `pms_brand`"
)

func (d *dao) RawBrandInfoList(ctx context.Context, ids []int64) (brands map[int64]*pb.Brand, err error) {
	if len(ids) == 0 {
		return
	}
	var idStr string
	for _, id := range ids {
		if len(idStr) != 0 {
			idStr += ","
		}
		idStr += strconv.FormatInt(id, 10)
	}
	sql := fmt.Sprintf(_brandInfoList, idStr)
	rows, err := d.db.Query(ctx, sql)
	if err != nil {
		log.Errorv(ctx, log.KV("func", "BrandList"), log.KV("event", "mysql_query"), log.KV("error", err), log.KV("sql", sql))
		return
	}
	defer rows.Close()
	brands = make(map[int64]*pb.Brand)
	for rows.Next() {
		brand := new(pb.Brand)
		if err = rows.Scan(
			&brand.Id,
			&brand.Name,
			&brand.FirstLetter,
			&brand.Sort,
			&brand.FactoryStatus,
			&brand.ShowStatus,
			&brand.ProductCount,
			&brand.ProductCommentCount,
			&brand.Logo,
			&brand.BigPic,
			&brand.BrandStory); err != nil {
			log.Errorv(ctx, log.KV("func", "BrandList"), log.KV("event", "mysql_scan"), log.KV("error", err), log.KV("sql", sql))
			return
		}
		brands[brand.Id] = brand
	}
	log.Infov(ctx, log.KV("event", "mysql_query"), log.KV("row_num", len(brands)), log.KV("sql", sql))
	return
}

func (d *dao) BrandList(ctx context.Context) (ids []int64, err error) {
	sql := fmt.Sprintf(_brandList)
	rows, err := d.db.Query(ctx, sql)
	if err != nil {
		log.Errorv(ctx, log.KV("event", "mysql_query"), log.KV("error", err), log.KV("sql", sql))
		return
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		if err = rows.Scan(&id); err != nil {
			log.Errorv(ctx, log.KV("event", "mysql_scan"), log.KV("error", err), log.KV("sql", sql))
			return
		}
		ids = append(ids, id)
	}
	log.Infov(ctx, log.KV("event", "mysql_query"), log.KV("row_num", len(ids)), log.KV("sql", sql))
	return
}
```

## 定义memcache缓存接口

```go
type _mc interface {
	// mc: -key=keyBrand -type=get
	CacheBrandInfoList(c context.Context, ids []int64) (map[int64]*pb.Brand, error)
	// mc: -key=keyBrand -expire=d.demoExpire
	AddCacheBrandInfoList(c context.Context, brands map[int64]*pb.Brand) (err error)
	// mc: -key=keyBrand
	DeleteBrandInfoListCache(c context.Context, ids []int64) (err error)
}

func keyBrand(id int64) string {
	return fmt.Sprintf("brand_%d", id)
}
```

## go generate