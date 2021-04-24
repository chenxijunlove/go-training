

## 作业：

我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层。为什么，应该怎么做请写出代码？

## 答：

我认为不应该Wrap sql.ErrNoRows，sql.ErrNoRows对上层业务来说，有可能并非是一个error，上层业务可能会吞掉这个错误，继续处理其他逻辑。且Warp后会加堆栈信息，且影响上层等值判断。



service 伪代码

```go
//GetUser 获取用户详情
func (s *service)GetUser(ctx context.Context,id int)(user *models.User,err error){
	//获取用户资料

	//获取用户积分
	points,err :=s.dao.GetUserPoints(ctx,id)
	if err != nil {
		if err != sql.ErrNoRows{
			return nil, errors.Wrap(err,"GetUserPoints error")
		}
	}
	//后续获取其他信息
	return
}
```



dao 伪代码

```go
//GetUserPoints 获取用户积分
func (d *dao)GetUserPoints(ctx context.Context,id int)(up *models.UserPoints,err error){
	up = &models.UserPoints{}
	err = d.db.Where("id = ?", id).First(up).Error
	return 
}
```

