package storage

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

func BenchmarkFindOriginalURL(b *testing.B) {
	mapStore, _ := NewFileStore("")
	arrayStore, _ := NewFileArrayStore("")
	uid := uuid.New()
	urls := map[string]string{
		"abcdef": "https://github.com/serjyuriev",
		"fedcba": "https://gitlab.com/servady",
		"lkasdj": "https://yandex.ru",
		"aslkqs": "https://google.com",
		"cpsoks": "https://vk.com",
		"qwrkml": "https://habr.com",
		"sdfkbj": "https://discord.com",
		"qlwknf": "https://gmail.com",
		"zkljns": "https://twitch.tv",
		"qkwnmd": "https://vscode.dev",
	}

	mapStore.InsertManyURLs(context.Background(), uid, urls)
	arrayStore.InsertManyURLs(context.Background(), uid, urls)

	b.ReportAllocs()
	b.ResetTimer()

	b.Run("map store", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			mapStore.FindOriginalURL(context.Background(), "qlwknf")
		}
	})

	b.Run("array store", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			arrayStore.FindOriginalURL(context.Background(), "qlwknf")
		}
	})
}

func BenchmarkFindByOriginalURL(b *testing.B) {
	mapStore, _ := NewFileStore("")
	arrayStore, _ := NewFileArrayStore("")
	uid := uuid.New()
	urls := map[string]string{
		"abcdef": "https://github.com/serjyuriev",
		"fedcba": "https://gitlab.com/servady",
		"lkasdj": "https://yandex.ru",
		"aslkqs": "https://google.com",
		"cpsoks": "https://vk.com",
		"qwrkml": "https://habr.com",
		"sdfkbj": "https://discord.com",
		"qlwknf": "https://gmail.com",
		"zkljns": "https://twitch.tv",
		"qkwnmd": "https://vscode.dev",
	}

	mapStore.InsertManyURLs(context.Background(), uid, urls)
	arrayStore.InsertManyURLs(context.Background(), uid, urls)

	b.ReportAllocs()
	b.ResetTimer()

	b.Run("map store", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			mapStore.FindByOriginalURL(context.Background(), "https://twitch.tv")
		}
	})

	b.Run("array store", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			arrayStore.FindByOriginalURL(context.Background(), "https://twitch.tv")
		}
	})
}

func BenchmarkInsertManyURLs(b *testing.B) {
	mapStore, _ := NewFileStore("")
	arrayStore, _ := NewFileArrayStore("")
	uid := uuid.New()
	urls := map[string]string{
		"abcdef": "https://github.com/serjyuriev",
		"fedcba": "https://gitlab.com/servady",
		"lkasdj": "https://yandex.ru",
		"aslkqs": "https://google.com",
		"cpsoks": "https://vk.com",
		"qwrkml": "https://habr.com",
		"sdfkbj": "https://discord.com",
		"qlwknf": "https://gmail.com",
		"zkljns": "https://twitch.tv",
		"qkwnmd": "https://vscode.dev",
	}

	b.ResetTimer()

	b.Run("map store", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			mapStore.InsertManyURLs(context.Background(), uid, urls)
		}
	})

	b.Run("array store", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			arrayStore.InsertManyURLs(context.Background(), uid, urls)
		}
	})
}

func BenchmarkFindURLsByUser(b *testing.B) {
	mapStore, _ := NewFileStore("")
	arrayStore, _ := NewFileArrayStore("")
	uid := uuid.New()
	urls := map[string]string{
		"abcdef": "https://github.com/serjyuriev",
		"fedcba": "https://gitlab.com/servady",
		"lkasdj": "https://yandex.ru",
		"aslkqs": "https://google.com",
		"cpsoks": "https://vk.com",
		"qwrkml": "https://habr.com",
		"sdfkbj": "https://discord.com",
		"qlwknf": "https://gmail.com",
		"zkljns": "https://twitch.tv",
		"qkwnmd": "https://vscode.dev",
	}
	mapStore.InsertManyURLs(context.Background(), uid, urls)
	arrayStore.InsertManyURLs(context.Background(), uid, urls)

	uid2 := uuid.New()
	urls2 := map[string]string{
		"akjsnd": "https://github.com/serjyuriev",
		"alkjnq": "https://gitlab.com/servady",
		"alskdm": "https://yandex.ru",
		"mklqwm": "https://google.com",
		"lvmzsq": "https://vk.com",
		"qwljeb": "https://habr.com",
		"kcqnaj": "https://discord.com",
		"klqnfk": "https://gmail.com",
		"mjklqc": "https://twitch.tv",
		"clkqns": "https://vscode.dev",
	}

	mapStore.InsertManyURLs(context.Background(), uid2, urls2)
	arrayStore.InsertManyURLs(context.Background(), uid2, urls2)

	b.ResetTimer()

	b.Run("map store", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			mapStore.FindURLsByUser(context.Background(), uid)
		}
	})

	b.Run("array store", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			arrayStore.FindURLsByUser(context.Background(), uid)
		}
	})
}
