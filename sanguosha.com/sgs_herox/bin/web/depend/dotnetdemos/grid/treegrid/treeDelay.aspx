<%@ Page Language="C#" AutoEventWireup="true" CodeFile="treeDelay.aspx.cs" Inherits="dotnetdemos_grid_treegrid_treeDelay" %>
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
    <title></title>
    <link href="../../../lib/ligerUI/skins/Aqua/css/ligerui-all.css" rel="stylesheet" type="text/css" />
 
    <script src="../../../lib/jquery/jquery-1.9.0.min.js" type="text/javascript"></script>
    <script src="../../../lib/ligerUI/js/core/base.js" type="text/javascript"></script>
    <script src="../../../lib/ligerUI/js/plugins/ligerGrid.js" type="text/javascript"></script>  

    <script type="text/javascript">



        var manager;
        $(function ()
        {
            window['g'] = 
            manager = $("#maingrid").ligerGrid({
                columns: [
                { display: '部门名', name: 'name', width: 150, align: 'left' },
                { display: '加入日期', name: 'date',format:"yyyy年MM月dd日", width: 150, type: 'date', align: 'left' },
                { display: '部门标示', name: 'id', width: 150, type: 'int', align: 'left' },
                { display: '部门描述', name: 'remark', width: 250, align: 'left' }
                ], width: '100%', pageSizeOptions: [5, 10, 15, 20], height: '97%',
                url: 'treeDelay.aspx?Action=GetData',
                dataAction: 'local',//本地排序
                usePager:false,
                alternatingRow: false, 
                tree: { columnName: 'name' },
                onTreeExpand: function (rowdata)
                {
                    console && console.log && console.log("onTreeExtend:" + rowdata['__id']);
                    var children = rowdata.children;
                    if (children.length == 0)
                    {
                        $.ajax({
                            url:  'treeDelay.aspx?Action=GetChildrenData' , data: {
                                pid:  rowdata.id
                            },
                            type: "POST"
                            , dataType: 'json'
                            , success: function (r)
                            {
                                for (var i = 0; i < r.Rows.length; i++)
                                {
                                    g.add(r.Rows[i], null, true, rowdata);
                                }
                                g.expand(rowdata);
                            }
                            , error: function ()
                            {
                                alert("error!");
                            }
                        });
                        return false;
                    }
                   
                },
                onTreeExpanded: function (rowdata)
                {
                    console && console.log && console.log("onTreeExpanded:" + rowdata['__id']);
                },
                onTreeCollapsed: function (rowdata)
                {
                    console && console.log && console.log("onTreeCollapsed:" + rowdata['__id']);
                },
                onTreeCollapse: function (rowdata)
                {
                    console && console.log && console.log("onTreeCollapse:" + rowdata['__id']);
                }
            }
            );
        });


        function getSelected()
        {
            var row = manager.getSelectedRow();
            if (!row) { alert('请选择行'); return; }
            alert(JSON.stringify(row));
        } 
        
    </script>
</head>
<body  style="padding:4px"> 
<div>  
  
   <a class="l-button" style="width:120px;float:left; margin-left:10px;" onclick="getSelected()">获取值</a>

   <div class="l-clear"></div>
 
</div>

    <div id="maingrid"></div>  
</body>
</html>