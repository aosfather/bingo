function showDialog(title,url) {
    layer.open({
        type: 2
        ,title: title
        ,area: ['480px', '560px']
        ,shade: 0
        ,maxmin: true
        ,content: url
        ,btn: []
        ,zIndex: layer.zIndex
        ,success: function(layero,index){
            console.log(index);
            console.log(layero);
        }
    });
}

function showDialogExt(title,url,obj){
    layer.open({
        type: 2
        ,title: title
        ,area: ['480px', '560px']
        ,shade: 0
        ,maxmin: true
        ,content: url
        ,btn: []
        ,zIndex: layer.zIndex //重点1
        ,success: function(layero,index){
            var body = layer.getChildFrame('body', index);
            var iframeWin = window[layero.find('iframe')[0]['name']]; //得到iframe页的窗口对象，执行iframe页的方法：iframeWin.method();
            console.log(body.html()) //得到iframe页的body内容
            body.find("input[name='id']").val('Hi，我是从父页来的')
            body.find("input").each(function(){this.value=obj[this.name];});
            console.log(index);
            console.log(layero);
            console.log(obj);
        }
    });
}