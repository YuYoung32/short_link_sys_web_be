# database 初始化指南

## 首次运行时

1. 创建数据库

   手动创建数据库`short_link_sys`, 和config.yaml里面的数据库名保持一致。可使用`sqls/create_database.sql`运行创建。

2. 创建表

   方法1：运行项目`short_link_sys_web`，会自动创建表`links`；运行项目`short_link_sys_core`会自动创建表`visits`。

   方法2：使用`sqls/create_tables.sql`创建。

3. 创建视图`link_visit_view`

   使用`sqls/create_link_visit_view.sql`创建。

4. 生成测试数据

   * 生成短链数据到表`links`

     运行`data_test.go`的`TestGenerateLinkData`

   * 生成访问数据到表`visits`

     运行`data_test.go`的`TestGenerateVisitData`
