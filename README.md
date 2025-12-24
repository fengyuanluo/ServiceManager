# Go Backend Template

这是一个单用户Go后端脚手架，用于快速构建Go后端项目。

## 集成组件

- [Gin](https://github.com/gin-gonic/gin)
- [Viper](https://github.com/spf13/viper)
- [Cobra](https://github.com/spf13/cobra)
- [Zap](https://github.com/uber-go/zap)
- [SQLite](https://github.com/go-gorm/sqlite)
- [JWT](https://github.com/golang-jwt/jwt)
- [CORS](https://github.com/gin-contrib/cors)
- [Swaggo](https://github.com/swaggo/swag)
- Static 
- Router
- Shutdown

## Swagger 文档

生成文档：

```bash
go generate ./...
```

访问地址：

- Swagger UI：`http://<host>:<port>/api/swagger/index.html`
- Swagger JSON：`http://<host>:<port>/api/swagger/doc.json`
