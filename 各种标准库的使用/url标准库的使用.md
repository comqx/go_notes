[TOC]

# 简介

# 使用

##QueryUnescape解码

```GO
abcd := `dimensions=%7BuserId%3D1513409009336235%2C+instanceId%3Di-j6c7s9rn1cj778v2jq3i%7D`
a, _ := url.QueryUnescape(abcd)
fmt.Println(a)
// dimensions={userId=1513409009336235, instanceId=i-j6c7s9rn1cj778v2jq3i}
```

# 注意

