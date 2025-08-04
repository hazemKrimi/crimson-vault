# Crimson Vault

Simple invoicing solution for small businesses that is open source and self-hosted written in Go.

# Development

This project uses [Air](https://github.com/air-verse/air) for live reloading so you can simply run the following command to get reloading on save:

```
air
```

For migrations the project is using [Atlas GORM Provider](https://github.com/ariga/atlas-provider-gorm) for migrations. To save a migration run the follwing command:

```
atlas migrate diff --env gorm
```

And to apply a migration on your development database you will need to run the following command:

```
atlas schema apply --env gorm -u "your database connection string"
```
