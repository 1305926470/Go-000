## Week02

如过 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层。
为什么，应该怎么做请写出代码？

方式一：dao 返回错误 `ErrNoRows` 错误, 在 biz 中 通过 `errors.Is(err, dao.ErrNoRows)`做判断，并响应的降级处理，如返回默认值。

方式二：通过  opaque error 的方式，对上层暴露一个 `IsErrNoRows() bool` 方法，用来判断是否为 `ErrNoRows`错误。不过这样的话，需要在biz 中先进行一次断言。

第一种方式需要对上层暴露一个 sentinel error ，增大了依赖面积；第二种方式需要进行断言处理。感觉都不够好。

其他类型的 db error 直接在 dao 中通过 errors.Wrap() 进行包装，然后直接往上透传，在上层拦截器（中间件）统一做错误日志处理。

