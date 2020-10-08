--[==[
  常用功能的lua脚本函数的封装
  1、result结果集处理


--]==]

function toSuccessResult(data)
  result={}
  result.code=0
  result.data={}
  for k,v in pairs(data) do
    table.insert(result.data, v)
  end
  return result
end

function toFailResult(code,msg)
  result={}
  result.code=code
  result.msg=msg
  return result
end