<script setup lang="ts">


import {computed,  onMounted, reactive, ref} from "vue";

import {EventsOn, HandleFormSubmission, ReadDirFileCount} from '../../wailsjs'
import {Message} from '@arco-design/web-vue';

const state = reactive({
  count: 0
})
const showProcess = ref(false)
const submitDisabled = ref(false)

function getFileCount() {
  ReadDirFileCount(form.sourceDir, form.includeChild).then(result => {
    state.count = result
  })
}

const form = reactive({
  sourceDir: ``,
  includeChild: false,
  targetDir: ``,
  dirType: '/20060102/', // \20200301\
  isMove: false, // false copy true move
  noShotTimeType: 1,//0 跳过 1 根据修改时间整理 2 根据创建时间整理
  skipSameFile: false, // 自动对比文件 不重复的重新命名
  skipFileLessThan: 10,//kb
  skipFileContains: '',//忽略名称
})
const rules = {
  sourceDir: [
    {
      required: true,
      message: 'sourceDir is required',
    },
  ],
  targetDir: [
    {
      required: true,
      message: 'targetDir is required',
    },
  ],
  dirType: [
    {
      required: true,
      message: 'dirType is required',
    },
  ],
}
const handleSubmit = ({values, errors}) => {
  console.log('values:', values, '\nerrors:', errors)

  if (errors != null) {

    return;
  }

  showProcess.value = true;
  submitDisabled.value = true

  Message.success('已启动照片整理!');

  console.log('表单提交:', form);
  // 在这里可以添加处理表单数据的逻辑

  HandleFormSubmission({
    ...form
  }).then(res => {

  });

};


const process = ref(0);
const total = ref(0);


const percent = computed(() => {

  if (total.value == 0 || process.value == 0) {
    return 0;
  } else {
    return process.value / total.value ?? 0;
  }
});

const message = ref('');

onMounted(() => {
  // 监听任务进度更新
  EventsOn('task-progress', (_process, _total) => {
    process.value = _process;
    total.value = _total;
    // Message.info('监听任务进度更新! ' + (percent.value * 100) + '%');
  });

  // 监听任务完成事件
  EventsOn('task-complete', (data) => {

    Message.success('照片整理完成!');

    showProcess.value = false;
    submitDisabled.value = false;
  });
});

// const layout = ref('vertical')
</script>

<template>
  <!--labelCol="{ span: 8 }" wrapperCol="{ span: 14 }" -->
  <a-form :model="form" :rules="rules" layout="vertical" style="width: 500px" @submit="handleSubmit" size="small">
    <a-form-item label="源文件夹" field="sourceDir" validate-trigger="input" required>
      <a-input v-model="form.sourceDir" placeholder="请输入源文件夹" @blur="getFileCount()"/>
      <template #extra>
        <div> 当前目录下共有 {{ state.count }} 个文件</div>
      </template>

    </a-form-item>

    <a-form-item label="包含子文件夹">
      <a-switch v-model="form.includeChild"/>
    </a-form-item>

    <a-form-item label="目标文件夹" field="targetDir" validate-trigger="input" required>
      <a-input v-model="form.targetDir" placeholder="请输入目标文件夹"/>
    </a-form-item>

    <a-form-item label="整理后的文件目录结构 (以2020年3月1日为例)" field="dirType" validate-trigger="blur" required>
      <a-select v-model="form.dirType" placeholder="请选择文件夹类型，如 \20200301\">

<!--golang的时间格式化需要2006年-->
        <a-option value="/2006/01/">\2020\03\</a-option>
        <a-option value="/2006/1/">\2020\3\</a-option>
        <a-option value="/2006/200601/">\2020\202003\</a-option>
        <a-option value="/200601/" style="border-top-style: solid;border-top-width: 1px">\202003\</a-option>
        <a-option value="/2006/01/02/" style="border-top-style: solid;border-top-width: 1px">\2020\03\01\</a-option>
        <a-option value="/2006/0102/">\2020\0301\</a-option>
        <a-option value="/20060102/">\20200301\</a-option>
        <a-option value="/2006/" style="border-top-style: solid;border-top-width: 1px">\2020\</a-option>

      </a-select>
    </a-form-item>

    <a-form-item label="整理文件的方法?">
      <a-radio-group v-model="form.isMove">
        <a-radio :value="false">复制</a-radio>
        <a-radio :value="true">移动</a-radio>
      </a-radio-group>
    </a-form-item>

    <a-form-item label="文件没有拍摄时间时?">
      <a-radio-group v-model="form.noShotTimeType">
        <a-radio :value="0">跳过</a-radio>
        <a-radio :value="1">根据修改时间整理</a-radio>
        <a-radio :value="2">根据创建时间整理</a-radio>
      </a-radio-group>
    </a-form-item>

    <a-form-item label="目标文件夹存在同名文件时?">
      <a-radio-group v-model="form.skipSameFile">
        <a-radio :value="false">智能处理(推荐)</a-radio>
        <a-radio :value="true">跳过</a-radio>
      </a-radio-group>
    </a-form-item>

    <a-form-item label="跳过小于文件大小">
      <a-input-number v-model="form.skipFileLessThan" :min="0"/>
      KB
    </a-form-item>

    <a-form-item label="忽略名称">
      <a-input v-model="form.skipFileContains" placeholder="请输入要忽略的文件名称"/>
    </a-form-item>

    <a-form-item>

      <a-space>
        <a-button html-type="submit" :disabled="submitDisabled">启动整理</a-button>
        <!--        <a-button @click="$refs.formRef.resetFields()">Reset</a-button>-->

      </a-space>
    </a-form-item>
    <a-form-item>
      <a-progress :percent="percent" :style="{width:'470px'}">
        <template v-slot:text="scope">
          {{ process }} / {{ total }} 进度 {{ (scope.percent * 100).toFixed(2) }}%
        </template>
      </a-progress>
    </a-form-item>
  </a-form>
</template>

<style scoped>

</style>
