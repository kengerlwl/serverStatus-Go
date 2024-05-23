<template>
    <a-layout class="layout">
      <a-layout-header>
        <div class="logo" />
        <a-menu
          v-model:selectedKeys="selectedKeys"
          theme="dark"
          mode="horizontal"
          :style="{ lineHeight: '64px' }"
        >
          <a-menu-item key="1">首页</a-menu-item>
          <a-menu-item key="2">服务</a-menu-item>
          <a-menu-item key="3">其他</a-menu-item>
        </a-menu>
      </a-layout-header>
      <a-layout-content style="padding: 0 50px">
        <a-breadcrumb style="margin: 16px 0">
          <!-- <a-breadcrumb-item>Home</a-breadcrumb-item>
          <a-breadcrumb-item>List</a-breadcrumb-item>
          <a-breadcrumb-item>App</a-breadcrumb-item> -->
        </a-breadcrumb>

    <!-- 元素居中 -->


        <div :style="{ background: '#fff', padding: '24px', minHeight: '280px', display: 'flex', justifyContent: 'center', alignItems: 'center' }">
            <serversCard />
        </div>

      </a-layout-content>



      <a-layout-footer style="text-align: center">
        A server probe project design by Kenger
      </a-layout-footer>
    </a-layout>
  </template>

<style>


html,
body,
#app{
  min-height: 100%;
  min-width: 100%;
  margin: 0;
  padding: 0;
}
</style>

<script setup>
import { ref } from 'vue';
import serversCard from './serversCard.vue';

const selectedKeys = ref(['2']);
</script>
<style scoped>
.site-layout-content {
min-height: 280px;
padding: 24px;
background: #fff;
min-width: 100%;
}
#components-layout-demo-top .logo {
float: left;
width: 100%;
height: 31px;
margin: 16px 24px 16px 0;
background: rgba(255, 255, 255, 0.3);
}
.ant-row-rtl #components-layout-demo-top .logo {
float: right;
margin: 16px 0 16px 24px;
}

[data-theme='dark'] .site-layout-content {
background: #141414;
}
</style>
