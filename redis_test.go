package service

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func TestRedis_Set(t *testing.T) {

	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	type fields struct {
		client *redis.Client
	}
	type args struct {
		key   string
		value interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "success",
			fields: fields{},
			args: args{
				key:   "test",
				value: "test value",
			},
			wantErr: false,
		},
	}

	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rds := &Redis{
				client: client,
			}
			if err := rds.Set(tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Redis.Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRedis_GetSting(t *testing.T) {

	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	type fields struct {
		client *redis.Client
	}
	type args struct {
		key          string
		defaultValue string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name:   "success",
			fields: fields{},
			args: args{
				key:          "test",
				defaultValue: "",
			},
			want:    "test value",
			wantErr: false,
		},
	}

	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	_ = client.Set(context.TODO(), "test", "test value", 0)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rds := &Redis{
				client: client,
			}
			got, err := rds.GetSting(tt.args.key, tt.args.defaultValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("Redis.GetSting() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Redis.GetSting() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedis_GetInt(t *testing.T) {

	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	type fields struct {
		client *redis.Client
	}
	type args struct {
		key          string
		defaultValue int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		{
			name:   "success",
			fields: fields{},
			args: args{
				key:          "test",
				defaultValue: 0,
			},
			want:    123,
			wantErr: false,
		},
	}

	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	_ = client.Set(context.TODO(), "test", "123", 0)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rds := &Redis{
				client: client,
			}
			got, err := rds.GetInt(tt.args.key, tt.args.defaultValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("Redis.GetInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Redis.GetInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedis_GetInt64(t *testing.T) {

	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	type fields struct {
		client *redis.Client
	}
	type args struct {
		key          string
		defaultValue int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		{
			name:   "success",
			fields: fields{},
			args: args{
				key:          "test",
				defaultValue: 0,
			},
			want:    12345678910,
			wantErr: false,
		},
	}

	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	_ = client.Set(context.TODO(), "test", "12345678910", 0)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rds := &Redis{
				client: client,
			}
			got, err := rds.GetInt64(tt.args.key, tt.args.defaultValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("Redis.GetInt64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Redis.GetInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedis_GetFloat64(t *testing.T) {

	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	type fields struct {
		client *redis.Client
	}
	type args struct {
		key          string
		defaultValue float64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    float64
		wantErr bool
	}{
		{
			name:   "success",
			fields: fields{},
			args: args{
				key:          "test",
				defaultValue: 0,
			},
			want:    12345678910,
			wantErr: false,
		},
	}

	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	_ = client.Set(context.TODO(), "test", "12345678910", 0)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rds := &Redis{
				client: client,
			}
			got, err := rds.GetFloat64(tt.args.key, tt.args.defaultValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("Redis.GetFloat64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Redis.GetFloat64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedis_HMGet(t *testing.T) {

	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	param := &RedisParameter{
		Network:      "tcp",
		Address:      s.Addr(),
		Password:     "",
		DB:           0,
		DialTimeout:  time.Duration(time.Second * 5),
		ReadTimeout:  time.Duration(time.Second * 5),
		WriteTimeout: time.Duration(time.Second * 5),
		PoolSize:     10,
	}
	rds, _ := NewRedis(param)

	key := "key"
	data := map[string]interface{}{
		"k1": "1",
		"k2": "2",
	}
	rds.HMSet(key, data)

	var fields []string

	for k := range data {
		fields = append(fields, k)
	}

	type args struct {
		key    string
		fields []string
	}

	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		{
			"0",
			args{key: key, fields: fields},
			data,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := rds.HMGet(tt.args.key, tt.args.fields)
			if (err != nil) != tt.wantErr {
				t.Errorf("Redis.HMGet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Redis.HMGet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedis_HMGetByFields(t *testing.T) {

	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	param := &RedisParameter{
		Network:      "tcp",
		Address:      s.Addr(),
		Password:     "",
		DB:           0,
		DialTimeout:  time.Duration(time.Second * 5),
		ReadTimeout:  time.Duration(time.Second * 5),
		WriteTimeout: time.Duration(time.Second * 5),
		PoolSize:     10,
	}
	rds, _ := NewRedis(param)

	key := "key"
	data := map[string]interface{}{
		"k1": "1", //redis 出來都是字串
		"k2": "2",
	}
	rds.HMSet(key, data)

	type args struct {
		key string
	}

	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		{
			"0",
			args{key: key},
			data,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := rds.HMGetByFields(tt.args.key, "k1", "k2")
			if (err != nil) != tt.wantErr {
				t.Errorf("Redis.HMGetByField() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Redis.HMGetByField() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedis_HExist_HGet(t *testing.T) {

	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	param := &RedisParameter{
		Network:      "tcp",
		Address:      s.Addr(),
		Password:     "",
		DB:           0,
		DialTimeout:  time.Duration(time.Second * 5),
		ReadTimeout:  time.Duration(time.Second * 5),
		WriteTimeout: time.Duration(time.Second * 5),
		PoolSize:     10,
	}
	rds, err := NewRedis(param)

	if err != nil {
		t.Errorf("NewRedis error = %v", err)
		return
	}

	key := "key"
	data := map[string]interface{}{
		"k1": 1,
		"k2": 2,
	}
	rds.HMSet(key, data)

	var fields []string

	for k := range data {
		fields = append(fields, k)
	}

	type args struct {
		key    string
		fields string
	}

	tests := []struct {
		name      string
		args      args
		want      string
		wantErr   bool
		wantExist bool
	}{
		{
			"0",
			args{key: key, fields: "no"},
			"1",
			false,
			false,
		},
		{
			"1",
			args{key: key, fields: "k1"},
			"1",
			false,
			true,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			fieldExist, err := rds.GetClient().HExists(context.TODO(), tt.args.key, tt.args.fields).Result()
			if err != nil {
				t.Fatalf("rds.GetClient().HExists() error = %v", err)
			}

			if fieldExist != tt.wantExist {
				t.Fatalf("rds.GetClient().HExist() exist = %v, wantExist %v", fieldExist, tt.wantExist)
			}

			if fieldExist {
				got, err := rds.GetClient().HGet(context.TODO(), tt.args.key, tt.args.fields).Result()

				if (err != nil) != tt.wantErr {
					t.Errorf("rds.GetClient().HGet() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("rds.GetClient().HGet() = %v, want %v", got, tt.want)
				}
			}
		})

	}
}

func TestRedis_Expire(t *testing.T) {

	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	rds := &Redis{
		client: client,
	}
	var key = "key1"

	//set key
	_ = rds.Set(key, 1)
	rds.client.Expire(context.TODO(), key, time.Duration(time.Second*10))

	str, _ := rds.client.Get(context.TODO(), key).Result()

	t.Logf("redis set variable be transfer to string %s", reflect.TypeOf(str))

	//ttl
	//會被set, getset 清除
	du, _ := rds.client.TTL(context.TODO(), key).Result()
	if du.Seconds() < 0 {
		t.Logf("redis key ttl (time to live) %f", du.Seconds())
	}

	//persist
	//ttl 會變成未設置，回 -1
	rds.client.Persist(context.TODO(), key)
	du, _ = rds.client.TTL(context.TODO(), key).Result()
	if du.Seconds() != -1 {
		t.Errorf("redis key after persist ttl (time to live) %f", du.Seconds())
	}

	bo, _ := rds.client.PExpire(context.TODO(), key, time.Duration(time.Millisecond*1)).Result()
	t.Logf("PExpire 1 millisecond %v", bo)

	k, err := rds.client.Get(context.TODO(), key).Result()
	t.Logf("key %s", k)
}

func TestRedis_Expire2(t *testing.T) {

	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	type fields struct {
		client *redis.Client
	}
	type args struct {
		key    string
		expire time.Duration
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name:   "key exist",
			fields: fields{},
			args: args{
				key:    "key1",
				expire: time.Duration((time.Second * 10)),
			},
			want: true,
		},
		{
			name:   "key doesn't exist",
			fields: fields{},
			args: args{
				key:    "key2",
				expire: time.Duration((time.Second * 10)),
			},
			want: false,
		},
	}

	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	_ = client.Set(context.TODO(), "key1", "123", 0)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rds := &Redis{
				client: client,
			}
			if got := rds.Expire(tt.args.key, tt.args.expire); got != tt.want {
				t.Errorf("Redis.Expire2() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedis_Exist(t *testing.T) {
	const Exists int64 = 1
	const DoesNotExists int64 = 0
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	type fields struct {
		client *redis.Client
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int64
	}{
		{
			name:   "key exist",
			fields: fields{},
			args: args{
				key: "key1",
			},
			want: Exists,
		},
		{
			name:   "key doesn't exist",
			fields: fields{},
			args: args{
				key: "key2",
			},
			want: DoesNotExists,
		},
	}

	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	_ = client.Set(context.TODO(), "key1", "123", 0)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rds := &Redis{
				client: client,
			}
			if got := rds.Exist(tt.args.key); got != tt.want {
				t.Errorf("Redis.Exist() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedis_Delete(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	type fields struct {
		client *redis.Client
	}
	type args struct {
		key []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		{
			name:   "delete 1 key",
			fields: fields{},
			args: args{
				key: []string{"key1"},
			},
			want:    1,
			wantErr: false,
		},
		{
			name:   "delete many keys",
			fields: fields{},
			args: args{
				key: []string{"key2", "key3"},
			},
			want:    2,
			wantErr: false,
		},
	}

	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	ctx := context.Background()
	_ = client.Set(ctx, "key1", "123", 0)
	_ = client.Set(ctx, "key2", "123", 0)
	_ = client.Set(ctx, "key3", "123", 0)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rds := &Redis{
				client: client,
			}
			got, err := rds.Delete(tt.args.key...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Redis.Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Redis.Delete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedis_HIncrBy(t *testing.T) {

	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	type fields struct {
		client *redis.Client
	}
	type args struct {
		key   string
		field string
		incr  int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		{
			name:   "key and field already exist",
			fields: fields{},
			args: args{
				key:   "key1",
				field: "node",
				incr:  1,
			},
			want:    2,
			wantErr: false,
		},
	}

	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	_ = client.HSet(context.TODO(), "key1", "node", 1)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rds := &Redis{
				client: client,
			}
			got, err := rds.HIncrBy(tt.args.key, tt.args.field, tt.args.incr)
			if (err != nil) != tt.wantErr {
				t.Errorf("Redis.HINCRBY() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Redis.HINCRBY() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedis_HGetSting(t *testing.T) {

	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})
	_ = client.HSet(context.TODO(), "key1", "field1", "string")

	type fields struct {
		client *redis.Client
	}
	type args struct {
		key          string
		field        string
		defaultValue string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name:   "success",
			fields: fields{client},
			args: args{
				key:          "key1",
				field:        "field1",
				defaultValue: "",
			},
			want:    "string",
			wantErr: false,
		},
		{
			name:   "failed and return default value",
			fields: fields{client},
			args: args{
				key:          "key1",
				field:        "field2",
				defaultValue: "",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rds := &Redis{
				client: tt.fields.client,
			}
			got, err := rds.HGetSting(tt.args.key, tt.args.field, tt.args.defaultValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("Redis.HGetSting() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Redis.HGetSting() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedis_HGetInt(t *testing.T) {

	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})
	_ = client.HSet(context.TODO(), "key1", "field1", 1)

	type fields struct {
		client *redis.Client
	}
	type args struct {
		key          string
		field        string
		defaultValue int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		{
			name:   "success",
			fields: fields{client},
			args: args{
				key:          "key1",
				field:        "field1",
				defaultValue: 0,
			},
			want:    1,
			wantErr: false,
		},
		{
			name:   "failed and return default value",
			fields: fields{client},
			args: args{
				key:          "key1",
				field:        "field2",
				defaultValue: 0,
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rds := &Redis{
				client: tt.fields.client,
			}
			got, err := rds.HGetInt(tt.args.key, tt.args.field, tt.args.defaultValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("Redis.HGetInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Redis.HGetInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedis_HGetInt64(t *testing.T) {

	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})
	_ = client.HSet(context.TODO(), "key1", "field1", 12345678910)

	type fields struct {
		client *redis.Client
	}
	type args struct {
		key          string
		field        string
		defaultValue int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		{
			name:   "success",
			fields: fields{client},
			args: args{
				key:          "key1",
				field:        "field1",
				defaultValue: 0,
			},
			want:    12345678910,
			wantErr: false,
		},
		{
			name:   "failed and return default value",
			fields: fields{client},
			args: args{
				key:          "key1",
				field:        "field2",
				defaultValue: 0,
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rds := &Redis{
				client: tt.fields.client,
			}
			got, err := rds.HGetInt64(tt.args.key, tt.args.field, tt.args.defaultValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("Redis.HGetInt64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Redis.HGetInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedis_HGetFloat64(t *testing.T) {

	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})
	_ = client.HSet(context.TODO(), "key1", "field1", 12.34)

	type fields struct {
		client *redis.Client
	}
	type args struct {
		key          string
		field        string
		defaultValue float64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    float64
		wantErr bool
	}{
		{
			name:   "success",
			fields: fields{client},
			args: args{
				key:          "key1",
				field:        "field1",
				defaultValue: 0,
			},
			want:    12.34,
			wantErr: false,
		},
		{
			name:   "failed and return default value",
			fields: fields{client},
			args: args{
				key:          "key1",
				field:        "field2",
				defaultValue: 0,
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rds := &Redis{
				client: tt.fields.client,
			}
			got, err := rds.HGetFloat64(tt.args.key, tt.args.field, tt.args.defaultValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("Redis.HGetFloat64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Redis.HGetFloat64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedis_HExistAndGetString(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})
	_ = client.HSet(context.TODO(), "key1", "field1", 12.34)

	type fields struct {
		client *redis.Client
	}
	type args struct {
		key    string
		fields string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		want1   bool
		wantErr bool
	}{
		{
			name:   "fields exist",
			fields: fields{client},
			args: args{
				key:    "key1",
				fields: "field1",
			},
			want:    "12.34",
			want1:   true,
			wantErr: false,
		},
		{
			name:   "fields not exist",
			fields: fields{client},
			args: args{
				key:    "key1",
				fields: "field2",
			},
			want:    "",
			want1:   false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rds := &Redis{
				client: tt.fields.client,
			}
			got, got1, err := rds.HExistAndGetString(tt.args.key, tt.args.fields)
			if (err != nil) != tt.wantErr {
				t.Errorf("Redis.HExistAndGetString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Redis.HExistAndGetString() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Redis.HExistAndGetString() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestRedis_LRange(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})
	_ = client.LPush(context.TODO(), "key1", "field1", 12.34)
	type fields struct {
		client *redis.Client
	}
	type args struct {
		key          string
		start        int64
		stop         int64
		defaultValue []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{
		{
			name:   "success",
			fields: fields{client: client},
			args: args{
				key:          "key1",
				start:        0,
				stop:         -1,
				defaultValue: []string{},
			},
			want:    []string{"12.34", "field1"},
			wantErr: false,
		},
		{
			name:   "key doesn't exist",
			fields: fields{client: client},
			args: args{
				key:          "key2",
				start:        0,
				stop:         -1,
				defaultValue: []string{"key", "doesn't", "exist"},
			},
			want:    []string{"key", "doesn't", "exist"},
			wantErr: true,
		},
		{
			name:   "illegal index",
			fields: fields{client: client},
			args: args{
				key:          "key1",
				start:        20,
				stop:         0,
				defaultValue: []string{"illegal", "index"},
			},
			want:    []string{"illegal", "index"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rds := &Redis{
				client: tt.fields.client,
			}
			got, err := rds.LRange(tt.args.key, tt.args.start, tt.args.stop, tt.args.defaultValue)
			if (err != nil) != tt.wantErr {
				t.Errorf("Redis.LRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Redis.LRange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedis_LPush(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	var value []interface{}
	value = append(value, 1)
	value = append(value, "value1")
	value = append(value, 2)
	value = append(value, "value2")

	type fields struct {
		client *redis.Client
	}
	type args struct {
		key   string
		value []interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "success",
			fields: fields{client: client},
			args: args{
				key:   "key1",
				value: value,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rds := &Redis{
				client: tt.fields.client,
			}
			if err := rds.LPush(tt.args.key, tt.args.value...); (err != nil) != tt.wantErr {
				t.Errorf("Redis.LPush() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRedis_IncrBy(t *testing.T) {

	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	type fields struct {
		client *redis.Client
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		{
			name:   "success",
			fields: fields{client: client},
			args: args{
				key: "key1",
			},
			want:    1,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rds := &Redis{
				client: tt.fields.client,
			}
			got, err := rds.IncrBy(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Redis.IncrBy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Redis.IncrBy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedis_TTL(t *testing.T) {

	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	type fields struct {
		client *redis.Client
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    time.Duration
		wantErr bool
	}{
		{
			name:   "failed - expired, key not exist",
			fields: fields{client: client},
			args: args{
				key: "key1",
			},
			want:    time.Duration(-2) * time.Second,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rds := &Redis{
				client: tt.fields.client,
			}
			got, err := rds.TTL(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("TTL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("TTL() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedis_HMGetAll(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	param := &RedisParameter{
		Network:      "tcp",
		Address:      s.Addr(),
		Password:     "",
		DB:           0,
		DialTimeout:  time.Duration(time.Second * 5),
		ReadTimeout:  time.Duration(time.Second * 5),
		WriteTimeout: time.Duration(time.Second * 5),
		PoolSize:     10,
	}
	rds, err := NewRedis(param)

	if err != nil {
		t.Errorf("NewRedis error = %v", err)
		return
	}

	key := "key"
	data := map[string]interface{}{
		"k1": 1,
		"k2": 2,
	}

	result := map[string]string{
		"k1": "1",
		"k2": "2",
	}

	rds.HMSet(key, data)

	type fields struct {
		client *redis.Client
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]string
		wantErr bool
	}{
		{
			name:    "success",
			fields:  fields{client: client},
			args:    args{key: key},
			want:    result,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rds := &Redis{
				client: tt.fields.client,
			}
			got, err := rds.HMGetAll(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Redis.HMGetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Redis.HMGetAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedis_GeoAdd(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	type fields struct {
		client *redis.Client
	}
	type args struct {
		key         string
		geoLocation []*GeoLocation
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				client: client,
			},
			args: args{
				key: "location",
				geoLocation: []*GeoLocation{
					{
						Name:      "location1",
						Longitude: 121.475,
						Latitude:  31.223,
					},
					{
						Name:      "location2",
						Longitude: 121.476,
						Latitude:  31.224,
					},
				},
			},
			want:    2,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rds := &Redis{
				client: tt.fields.client,
			}
			got, err := rds.GeoAdd(tt.args.key, tt.args.geoLocation...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Redis.GeoAdd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Redis.GeoAdd() = %v, want %v", got, tt.want)
			}
		})
	}
}

//NOTE: miniredis doesn't support geo search
func TestRedis_GeoSearch(t *testing.T) {
	// s, err := miniredis.Run()
	// if err != nil {
	// 	panic(err)
	// }
	// defer s.Close()

	// client := redis.NewClient(&redis.Options{
	// 	Addr: s.Addr(),
	// })

	// location := []*redis.GeoLocation{
	// 	{
	// 		Name:      "location1",
	// 		Longitude: -122.27652,
	// 		Latitude:  37.805186,
	// 	},
	// 	{
	// 		Name:      "location2",
	// 		Longitude: -122.2674626,
	// 		Latitude:  37.8062344,
	// 	},
	// 	{
	// 		Name:      "you can't find me",
	// 		Longitude: 128.2674626,
	// 		Latitude:  50.8062344,
	// 	},
	// }

	// _ = client.GeoAdd(context.TODO(), "location", location...)

	// type fields struct {
	// 	client *redis.Client
	// }
	// type args struct {
	// 	key string
	// 	q   *GeoSearchQuery
	// }
	// tests := []struct {
	// 	name    string
	// 	fields  fields
	// 	args    args
	// 	want    []string
	// 	wantErr bool
	// }{
	// 	{
	// 		name: "success",
	// 		fields: fields{
	// 			client: client,
	// 		},
	// 		args: args{
	// 			key: "location",
	// 			q: &GeoSearchQuery{
	// 				Member:     "",
	// 				Longitude:  -122.2612767,
	// 				Latitude:   37.7936847,
	// 				Radius:     5,
	// 				RadiusUnit: "km",
	// 				BoxWidth:   0.0,
	// 				BoxHeight:  0.0,
	// 				BoxUnit:    "",
	// 				Sort:       "",
	// 				Count:      0,
	// 				CountAny:   false,
	// 			},
	// 		},
	// 		want:    []string{"location1", "location2"},
	// 		wantErr: false,
	// 	},
	// }
	// for _, tt := range tests {
	// 	t.Run(tt.name, func(t *testing.T) {
	// 		rds := &Redis{
	// 			client: tt.fields.client,
	// 		}
	// 		got, err := rds.GeoSearch(tt.args.key, tt.args.q)
	// 		if (err != nil) != tt.wantErr {
	// 			t.Errorf("Redis.GeoSearch() error = %v, wantErr %v", err, tt.wantErr)
	// 			return
	// 		}
	// 		if !reflect.DeepEqual(got, tt.want) {
	// 			t.Errorf("Redis.GeoSearch() = %v, want %v", got, tt.want)
	// 		}
	// 	})
	// }
}

func TestRedis_SAdd(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	type fields struct {
		client *redis.Client
	}
	type args struct {
		key     string
		members []any
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				client: client,
			},
			args: args{
				key:     "test",
				members: []any{"1", "2"},
			},
			want:    2,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rds := &Redis{
				client: tt.fields.client,
			}
			got, err := rds.SAdd(tt.args.key, tt.args.members...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Redis.SAdd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Redis.SAdd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedis_SMembers(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	client.SAdd(context.Background(), "test", "1", "2")

	type fields struct {
		client *redis.Client
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				client: client,
			},
			args: args{
				key: "test",
			},
			want:    []string{"1", "2"},
			wantErr: false,
		},
		{
			name: "key not exist",
			fields: fields{
				client: client,
			},
			args: args{
				key: "test2",
			},
			want:    []string{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rds := &Redis{
				client: tt.fields.client,
			}
			got, err := rds.SMembers(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Redis.SMembers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Redis.SMembers() = %v, want %v", got, tt.want)
			}
		})
	}
}
