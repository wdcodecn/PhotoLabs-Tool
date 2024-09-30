<script lang="ts" setup>
import PhotoOrganization from "./views/PhotoOrganization.vue";
</script>

<template>

  <a-layout>
    <a-layout-header  >
<!--      <Header/>-->
    </a-layout-header>
    <a-layout-content style="padding: 0 20px;">

      <a-tabs default-active-key="1" lazy-load type="rounded" size="small">
        <a-tab-pane key="1" title="照片整理">
          <PhotoOrganization/>
        </a-tab-pane>
        <a-tab-pane key="2" title="精确去重">
          Content of Tab Panel 2
        </a-tab-pane>
        <a-tab-pane key="3" title="帮助">
          Content of Tab Panel 3
        </a-tab-pane>
        <a-tab-pane key="4" title="打赏作者">
          Content of Tab Panel 4
        </a-tab-pane>
      </a-tabs>

    </a-layout-content>
    <a-layout-footer  >
<!--     <Footer/>  -->
    </a-layout-footer>
  </a-layout>

</template>
