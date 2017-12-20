package credentialstore

type Config struct {
    port int
    region string
}

func GetConfig() Config {
    return Config{ 8080, "us-west-2" }
}