package vue

import (
	"fmt"
	"os"
	"path"
	"text/template"

	"github.com/abulo/ratel/v3/toolkit/base"
	"github.com/abulo/ratel/v3/util"
	"github.com/fatih/color"
	"github.com/spf13/cast"
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
		"PropsSearch":           base.PropsSearch,
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
	if exists := base.Config.Exists("template.VuePage"); exists {
		filePath := path.Join(cast.ToString(base.Path), base.Config.String("template.VuePage"))
		if util.FileExists(filePath) {
			if tplString, err := util.FileGetContents(filePath); err == nil {
				return tplString
			}
		}
	}
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
      {{- if InMethod .Method "Show"}}
      <el-button v-auth="'{{.Pkg}}.{{CamelStr .Table.TableName}}'" type="primary" link :icon="View" @click="handleItem(scope.row)">
        查看
      </el-button>
      {{- end}}
      <el-dropdown trigger="click">
        <el-button
              v-auth="[
                {{- if InMethod .Method "Update"}}
                '{{.Pkg}}.{{CamelStr .Table.TableName}}Update',
                {{- end}}
                {{- if InMethod .Method "Delete"}}
                '{{.Pkg}}.{{CamelStr .Table.TableName}}Delete',
                {{- end}}
                {{- if InMethod .Method "Recover"}}
                '{{.Pkg}}.{{CamelStr .Table.TableName}}Recover',
                {{- end}}
                {{- if InMethod .Method "Drop"}}
                '{{.Pkg}}.{{CamelStr .Table.TableName}}Drop',
                {{- end}}
              ]"
              type="primary"
              link
              :icon="DArrowRight">
              更多
        </el-button>
        <template #dropdown>
          <el-dropdown-menu>
            {{- if InMethod .Method "Update"}}
            <el-dropdown-item v-auth="'{{.Pkg}}.{{CamelStr .Table.TableName}}Update'" :icon="EditPen" @click="handleUpdate(scope.row)">
              编辑
            </el-dropdown-item>
            {{- end}}
            {{- if InMethod .Method "Delete"}}
            <el-dropdown-item
              v-if="scope.row.deleted === 0"
              v-auth="'{{.Pkg}}.{{CamelStr .Table.TableName}}Delete'"
              :icon="Delete"
              @click="handleDelete(scope.row)">
              删除
            </el-dropdown-item>
            {{- end}}
            {{- if InMethod .Method "Recover"}}
            <el-dropdown-item
              v-if="scope.row.deleted === 1"
              v-auth="'{{.Pkg}}.{{CamelStr .Table.TableName}}Recover'"
              :icon="Refresh"
              @click="handleRecover(scope.row)">
              恢复
            </el-dropdown-item>
            {{- end}}
            {{- if InMethod .Method "Drop"}}
            <el-dropdown-item
              v-if="scope.row.deleted === 1"
              v-auth="'{{.Pkg}}.{{CamelStr .Table.TableName}}Drop'"
              :icon="DeleteFilled"
              @click="handleDrop(scope.row)">
              清理
            </el-dropdown-item>
            {{- end}}
          </el-dropdown-menu>
        </template>
      </el-dropdown>
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
          <el-input v-model="{{Helper $.Table.TableName}}ItemFrom.{{Helper .ColumnName}}" :disabled="disabled" />
        </el-form-item>
		{{- end}}
	  </el-form>
	  <template #footer v-if="!disabled">
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
	{{- if InMethod .Method "Drop"}}
	DeleteFilled,
	{{- end}}
	{{- if InMethod .Method "Show"}}
	View,
	{{- end}}
	DArrowRight,
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
  {{- if InMethod .Method "Drop"}}
  drop{{CamelStr .Table.TableName}}Api,
  {{- end}}
  {{- if InMethod .Method "Recover"}}
  recover{{CamelStr .Table.TableName}}Api,
  {{- end}}
  {{- if InMethod .Method "Show"}}
  get{{CamelStr .Table.TableName}}Api,
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
const disabled = ref(true);
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
        span: 2,
		    props: {
          activeValue: 1,
          inactiveValue: 0
        }
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
		isShow: HasPermission(
      {{- if InMethod .Method "Update"}}
      "{{.Pkg}}.{{CamelStr .Table.TableName}}Update",
      {{- end}}
      {{- if InMethod .Method "Delete"}}
      "{{.Pkg}}.{{CamelStr .Table.TableName}}Delete",
      {{- end}}
      {{- if InMethod .Method "Drop"}}
      "{{.Pkg}}.{{CamelStr .Table.TableName}}Drop",
      {{- end}}
      {{- if InMethod .Method "Recover"}}
      "{{.Pkg}}.{{CamelStr .Table.TableName}}Recover",
      {{- end}}
      {{- if InMethod .Method "Show"}}
      "{{.Pkg}}.{{CamelStr .Table.TableName}}",
      {{- end}}
    )}
];

// 重置数据
const reset = () => {
  loading.value = false;
  {{Helper .Table.TableName}}ItemFrom.value = {
    {{Json .TableColumn}}
  };
  disabled.value = true;
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

{{- if InMethod .Method "Drop"}}
// 清理按钮
const handleDrop = async (row: {{CamelStr .Table.TableName}}.Res{{CamelStr .Table.TableName}}Item) => {
  await useHandleData(drop{{CamelStr .Table.TableName}}Api, Number(row.id), "清理{{.Table.TableComment}}");
  proTable.value?.getTableList();
};
{{- end}}

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
  disabled.value = false;
};
{{- end}}

{{- if InMethod .Method "Update"}}
// 编辑按钮
const handleUpdate = async (row: {{CamelStr .Table.TableName}}.Res{{CamelStr .Table.TableName}}Item) => {
  title.value = "编辑{{.Table.TableComment}}";
  centerDialogVisible.value = true;
  reset();
  const { data } = await get{{CamelStr .Table.TableName}}Api(Number(row.id));
  {{Helper .Table.TableName}}ItemFrom.value = data;
  disabled.value = false;
};
{{- end}}

{{- if InMethod .Method "Show"}}
// 查看按钮
const handleItem = async (row: {{CamelStr .Table.TableName}}.Res{{CamelStr .Table.TableName}}Item) => {
  title.value = "查看{{.Table.TableComment}}";
  centerDialogVisible.value = true;
  reset();
  const { data } = await get{{CamelStr .Table.TableName}}Api(Number(row.id));
  {{Helper .Table.TableName}}ItemFrom.value = data;
  disabled.value = true;
};
{{- end}}

</script>
<style lang="scss">
@import "@/styles/custom.scss";
</style>`
	return outString
}
