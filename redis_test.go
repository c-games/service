package service

import (
	"reflect"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis"
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

	_ = client.Set("test", "test value", 0)

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

	_ = client.Set("test", "123", 0)

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

	_ = client.Set("test", "12345678910", 0)

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

	_ = client.Set("test", "12345678910", 0)

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

	param := &redisParameter{
		Network:      "tcp",
		Address:      s.Addr(),
		Password:     "",
		DB:           0,
		DialTimeout:  time.Duration(time.Second * 5),
		ReadTimeout:  time.Duration(time.Second * 5),
		WriteTimeout: time.Duration(time.Second * 5),
		PoolSize:     10,
	}
	rds, _ := newRedis(param)

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

	param := &redisParameter{
		Network:      "tcp",
		Address:      s.Addr(),
		Password:     "",
		DB:           0,
		DialTimeout:  time.Duration(time.Second * 5),
		ReadTimeout:  time.Duration(time.Second * 5),
		WriteTimeout: time.Duration(time.Second * 5),
		PoolSize:     10,
	}
	rds, _ := newRedis(param)

	key := "key"
	data := map[string]interface{}{
		"k1": 1,
		"k2": 2,
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
	rds.client.Expire(key, time.Duration(time.Second*10))

	str, _ := rds.client.Get(key).Result()

	t.Logf("redis set variable be transfer to string %s", reflect.TypeOf(str))

	//ttl
	//會被set, getset 清除
	du, _ := rds.client.TTL(key).Result()
	if du.Seconds() < 0 {
		t.Logf("redis key ttl (time to live) %f", du.Seconds())
	}

	//persist
	//ttl 會變成未設置，回 -1
	rds.client.Persist(key)
	du, _ = rds.client.TTL(key).Result()
	if du.Seconds() != -1 {
		t.Errorf("redis key after persist ttl (time to live) %f", du.Seconds())
	}

	bo, _ := rds.client.PExpire(key, time.Duration(time.Millisecond*1)).Result()
	t.Logf("PExpire 1 millisecond %v", bo)

	k, err := rds.client.Get(key).Result()
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

	_ = client.Set("key1", "123", 0)

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
				key: "key1",
			},
			want: DoesNotExists,
		},
	}

	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	_ = client.Set("key2", "123", 0)

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
			want: 1,
			wantErr: false,
		},
		{
			name:   "delete many keys",
			fields: fields{},
			args: args{
				key: []string{"key2", "key3"},
			},
			want: 2,
			wantErr: false,
		},
	}

	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	_ = client.Set("key1", "123", 0)
	_ = client.Set("key2", "123", 0)
	_ = client.Set("key3", "123", 0)

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
