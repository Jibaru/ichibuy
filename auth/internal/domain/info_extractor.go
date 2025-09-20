package domain

type InfoExtractor func(token string) (string, string, error)
