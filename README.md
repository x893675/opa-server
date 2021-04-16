## OPA Server

使用 Open Policy Agent 实现的 **RBAC** 鉴权服务

参考:

* [OPA docs](https://www.openpolicyagent.org/docs/latest/)
* [repo playground](https://play.openpolicyagent.org/)
* [opa rbac](https://github.com/ashutoshSce/opa-rbac)

## 本地调试运行

1. `docker-compose up -d` 启动 opa server
2. 在 postman 中导入 `opa.postman_collection.json`
3. 根据导入的请求访问 opa server (注意更换服务地址)
4. 或者执行 `opa test -v api.rego api_test.rego` 运行测试

## Roadmap

- [ ] 更新 README.md
- [ ] 使用 [push-data 方式](https://www.openpolicyagent.org/docs/latest/external-data/#option-4-push-data) 实现 opa server 的 policy 和 data 的更新

