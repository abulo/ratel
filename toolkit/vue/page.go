package vue

import (
	"fmt"
	"os"
	"path"
	"text/template"

	"github.com/abulo/ratel/v3/toolkit/base"
	"github.com/abulo/ratel/v3/util"
	"github.com/fatih/color"
)

// viewUrl: 页面路径
func GeneratePage(moduleParam base.ModuleParam, fullPageDir, viewUrl, tableName string) {
	// 文件夹路径
	dir := util.GetParentDirectory(path.Join(fullPageDir, viewUrl))
	// 创建文件夹
	_ = os.MkdirAll(dir, os.ModePerm)
	// 文件路径
	viewFullFile := path.Join(fullPageDir, viewUrl)

	// 模板变量
	tpl := template.Must(template.New("page").Funcs(template.FuncMap{
		"Convert":               base.Convert,
		"SymbolChar":            base.SymbolChar,
		"Char":                  base.Char,
		"Helper":                base.Helper,
		"CamelStr":              base.CamelStr,
		"Add":                   base.Add,
		"ModuleProtoConvertDao": base.ModuleProtoConvertDao,
		"ModuleDaoConvertProto": base.ModuleDaoConvertProto,
		"ModuleProtoConvertMap": base.ModuleProtoConvertMap,
		"ApiToProto":            base.ApiToProto,
		"TypeScriptCondition":   base.TypeScriptCondition,
		"TypeScript":            base.TypeScript,
		"Json":                  base.Json,
		"InMethod":              base.InMethod,
		"Rule":                  base.Rule,
		"Props":                 base.Props,
	}).Parse(PageTemplate()))

	// 文件夹路径
	if util.FileExists(viewFullFile) {
		util.Delete(viewFullFile)
	}

	file, err := os.OpenFile(viewFullFile, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		fmt.Println("文件句柄错误:", color.RedString(err.Error()))
		return
	}

	//渲染输出
	err = tpl.Execute(file, moduleParam)
	if err != nil {
		fmt.Println("模板解析错误:", color.RedString(err.Error()))
		return
	}
	fmt.Printf("\n🍺 CREATED   %s\n", color.GreenString(viewFullFile))
}

func PageTemplate() string {
	outString := `<template>
  <div class="table-box">
	<ProTable
	  ref="proTable"
	  title="{{.Table.TableComment}}列表"
	  row-key="id"
	  :columns="columns"
	  :request-api="get{{CamelStr .Table.TableName}}ListApi"
	  :request-auto="true"
	  :pagination="{{.Page}}"
	  :search-col="12">
	  <!-- 表格 header 按钮 -->
      <template #tableHeader>
	  	{{- if InMethod .Method "Create"}}
        <el-button v-auth="'{{.Pkg}}.{{CamelStr .Table.TableName}}Create'" type="primary" :icon="CirclePlus" @click="handleAdd">新增</el-button>
		{{- end}}
      </template>
	  {{- if .SoftDelete}}
	  <!-- 删除状态 -->
      <template #deleted="scope">
        <DictTag type="delete" :value="scope.row.deleted" />
      </template>
	  {{- end}}
	  <!-- 菜单操作 -->
	  <template #operation="scope">
	  	{{- if InMethod .Method "Update"}}
        <el-button v-auth="'{{.Pkg}}.{{CamelStr .Table.TableName}}Update'" type="primary" link :icon="EditPen" @click="handleUpdate(scope.row)">
          编辑
        </el-button>
		{{- end}}
		{{- if InMethod .Method "Delete"}}
        <el-button
          v-if="scope.row.deleted === 0"
          v-auth="'{{.Pkg}}.{{CamelStr .Table.TableName}}Delete'"
          type="primary"
          link
          :icon="Delete"
          @click="handleDelete(scope.row)">
          删除
        </el-button>
		{{- end}}
		{{- if InMethod .Method "Recover"}}
        <el-button
          v-if="scope.row.deleted === 1"
          v-auth="'{{.Pkg}}.{{CamelStr .Table.TableName}}Recover'"
          type="primary"
          link
          :icon="Refresh"
          @click="handleRecover(scope.row)">
          恢复
        </el-button>
		{{- end}}
      </template>
	</ProTable>
	<el-dialog
      v-model="centerDialogVisible"
      :title="title"
      width="40%"
      destroy-on-close
      align-center
      center
      append-to-body
      draggable
      :lock-scroll="false"
      class="dialog-settings">
	  <el-form ref="ref{{CamelStr .Table.TableName}}ItemFrom" :model="{{Helper .Table.TableName}}ItemFrom" :rules="rules{{CamelStr .Table.TableName}}ItemFrom" label-width="100px">
		{{- range .TableColumn}}
		<el-form-item label="{{.ColumnComment}}" prop="{{Helper .ColumnName}}">
          <el-input v-model="{{Helper $.Table.TableName}}ItemFrom.{{Helper .ColumnName}}" />
        </el-form-item>
		{{- end}}
	  </el-form>
	  <template #footer>
		<span class="dialog-footer">
		  <el-button @click="resetForm(ref{{CamelStr .Table.TableName}}ItemFrom)">取消</el-button>
		  <el-button type="primary" :loading="loading" @click="submitForm(ref{{CamelStr .Table.TableName}}ItemFrom)">确定</el-button>
		</span>
	  </template>
	</el-dialog>
  </div>
</template>
<script setup lang="ts" name="{{Helper .Table.TableName}}">
import { ref, reactive } from "vue";
import { ProTableInstance, ColumnProps, SearchProps } from "@/components/ProTable/interface";
import {
	{{- if InMethod .Method "Update"}}
	EditPen,
	{{- end}}
	{{- if InMethod .Method "Create"}}
	CirclePlus,
	{{- end}}
	{{- if InMethod .Method "Delete"}}
	Delete,
	{{- end}}
	{{- if InMethod .Method "Recover"}}
	Refresh,
	{{- end}}
	} from "@element-plus/icons-vue";
import ProTable from "@/components/ProTable/index.vue";
import { {{CamelStr .Table.TableName}} } from "@/api/interface/{{Helper .Table.TableName}}";
import {
  {{- if InMethod .Method "List"}}
  get{{CamelStr .Table.TableName}}ListApi,
  {{- end}}
  {{- if InMethod .Method "Delete"}}
  delete{{CamelStr .Table.TableName}}Api,
  {{- end}}
  {{- if InMethod .Method "Recover"}}
  recover{{CamelStr .Table.TableName}}Api,
  {{- end}}
  {{- if InMethod .Method "Show"}}
  get{{CamelStr .Table.TableName}}ItemApi,
  {{- end}}
  {{- if InMethod .Method "Create"}}
  add{{CamelStr .Table.TableName}}Api,
  {{- end}}
  {{- if InMethod .Method "Update"}}
  update{{CamelStr .Table.TableName}}Api,
  {{- end}}
} from "@/api/modules/{{Helper .Table.TableName}}";
import { FormInstance, FormRules } from "element-plus";
{{- if .SoftDelete}}
import { getIntDictOptions } from "@/utils/dict";
import { DictTag } from "@/components/DictTag";
{{- end}}
import { useHandleData, useHandleSet } from "@/hooks/useHandleData";
import { HasPermission } from "@/utils/permission";
//加载
const loading = ref(false);
//弹出层标题
const title = ref();
//table数据
const proTable = ref<ProTableInstance>();
//是否显示弹出层
const centerDialogVisible = ref(false);
//数据接口
const {{Helper .Table.TableName}}ItemFrom = ref<{{CamelStr .Table.TableName}}.Res{{CamelStr .Table.TableName}}Item>({
{{Json .TableColumn}}
});
//校验
const ref{{CamelStr .Table.TableName}}ItemFrom = ref<FormInstance>();
//校验
const rules{{CamelStr .Table.TableName}}ItemFrom = reactive<FormRules>({
  {{Rule .TableColumn}}
});

{{- if .SoftDelete}}
//删除状态
const deletedEnum = getIntDictOptions("delete");
// 表格配置项
const deleteSearch = reactive<SearchProps>(
  HasPermission("{{.Pkg}}.{{CamelStr .Table.TableName}}Delete")
    ? {
        el: "switch",
        span: 2
      }
    : {}
);
{{- end}}

const columns: ColumnProps<{{CamelStr .Table.TableName}}.Res{{CamelStr .Table.TableName}}Item>[] = [
	  {{Props .TableColumn .Condition}}
	  {
		prop: "operation",
		label: "操作",
		width: 150,
		fixed: "right",
		isShow: HasPermission("{{.Pkg}}.{{CamelStr .Table.TableName}}Update", "{{.Pkg}}.{{CamelStr .Table.TableName}}Delete", "{{.Pkg}}.{{CamelStr .Table.TableName}}Recover")
	  }
];

// 重置数据
const reset = () => {
  loading.value = false;
  {{Helper .Table.TableName}}ItemFrom.value = {
    {{Json .TableColumn}}
  };
};

// resetForm
const resetForm = (formEl: FormInstance | undefined) => {
  centerDialogVisible.value = false;
  if (!formEl) return;
  formEl.resetFields();
};


// 提交数据
const submitForm = (formEl: FormInstance | undefined) => {
  if (!formEl) return;
  formEl.validate(async valid => {
    if (!valid) return;
    loading.value = true;
    const data = {{Helper .Table.TableName}}ItemFrom.value as unknown as {{CamelStr .Table.TableName}}.Res{{CamelStr .Table.TableName}}Item;
    if (data.id !== 0) {
	  {{- if InMethod .Method "Update"}}
      await useHandleSet(update{{CamelStr .Table.TableName}}Api, data.id, data, "修改{{.Table.TableComment}}");
	  {{- end}}
    } else {
	  {{- if InMethod .Method "Create"}}
      await useHandleData(add{{CamelStr .Table.TableName}}Api, data, "添加{{.Table.TableComment}}");
	  {{- end}}
    }
    resetForm(formEl);
    loading.value = false;
    proTable.value?.getTableList();
  });
};

{{- if InMethod .Method "Delete"}}
// 删除按钮
const handleDelete = async (row: {{CamelStr .Table.TableName}}.Res{{CamelStr .Table.TableName}}Item) => {
  await useHandleData(delete{{CamelStr .Table.TableName}}Api, Number(row.id), "删除{{.Table.TableComment}}");
  proTable.value?.getTableList();
};
{{- end}}

{{- if InMethod .Method "Recover"}}
// 恢复按钮
const handleRecover = async (row: {{CamelStr .Table.TableName}}.Res{{CamelStr .Table.TableName}}Item) => {
  await useHandleData(recover{{CamelStr .Table.TableName}}Api, Number(row.id), "恢复{{.Table.TableComment}}");
  proTable.value?.getTableList();
};
{{- end}}

{{- if InMethod .Method "Create"}}
// 添加按钮
const handleAdd = () => {
  title.value = "新增{{.Table.TableComment}}";
  centerDialogVisible.value = true;
  reset();
};
{{- end}}

// 编辑按钮
const handleUpdate = async (row: {{CamelStr .Table.TableName}}.Res{{CamelStr .Table.TableName}}Item) => {
  title.value = "编辑{{.Table.TableComment}}";
  centerDialogVisible.value = true;
  reset();
  const { data } = await get{{CamelStr .Table.TableName}}ItemApi(Number(row.id));
  {{Helper .Table.TableName}}ItemFrom.value = data;
};
</script>
<style lang="scss">
@import "@/styles/custom.scss";
</style>`
	return outString
}
