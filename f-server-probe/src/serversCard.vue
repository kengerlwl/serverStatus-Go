<template>
  <div :style="{ background: '#fff', padding: '24px', minHeight: '280px', display: 'flex', flexWrap: 'wrap', justifyContent: 'center', alignItems: 'flex-start', gap: '16px' }">
    <div v-for="server in servers" :key="server.host" style="{ flex: '0 1 auto' }">
        <a-card :title="`Server: ${server.name}`" style="width: 300px;background:rgb(248, 251, 251);">
          <template #extra>
            <!-- <a-button type="primary" :visible="server.server_online" v-model:visible="server.server_online" danger ghost>Danger</a-button> -->
            <!-- <a-button danger shape="round" loading /> -->
            <a-tag color="error" :visible="!server.server_online">error</a-tag>

            <a @click="showModal(server)">详细</a>
          </template>
         


            <div>
                <a-space direction="vertical" style="width: auto; alignItems:center">

                    <p style="margin: 0%; display: flex; ">CPU</p>
                    <a-progress
                        type="circle"
                        :stroke-color="{
                            '0%': '#108ee9',
                            '100%': '#87d068',
                        }"
                        :percent="(server.cpu).toFixed(1)"
                        :format="percent => `${percent}%`"
                        :size="80"
                        />
                    </a-space>

                <a-space direction="vertical" style="width: auto; alignItems:center">
                    
                    <p style="margin: 0%; display: flex; ">GPU</p>
                    <a-progress
                        type="circle"
                        :stroke-color="{
                            '0%': '#108ee9',
                            '100%': '#87d068',
                        }"
                        :percent="(server.gpu_mem_info * 100).toFixed(1)"
                        :format="percent => `${percent}%`"
                        :size="80"
                        />
                    </a-space>

                <a-space direction="vertical" style="width: auto; alignItems:center">
                    
                    <p style="margin: 0%; display: flex; ">内存</p>

                    <a-progress
                    type="circle"
                    :stroke-color="{
                        '0%': '#108ee9',
                        '100%': '#87d068',
                    }"
                    :percent="(server.memory_used / server.memory_total * 100).toFixed(1)"
                    :format="percent => `${percent}%`"
                    :size="80"
                    />
                </a-space>


                <a-row>
                    <!-- 要保证20 + 4 = 24. 根据这个来分配比例 -->
                    <a-col :span="20" :push="4">
                    
                        <a-progress
                            :stroke-color="{
                                '0%': '#108ee9',
                                '100%': '#87d068',
                            }"
                            :style="{width: '100%'}"
                            :percent="(server.hdd_used / server.hdd_total * 100).toFixed(1)"
                            :format="percent => `${percent}%`"
                            />
                    </a-col>
                    <a-col :span="4" :pull="20">
                        HHD：
                    </a-col>
                </a-row>


                <a-row>
                    <a-col :span="12">
                    <a-statistic title="上传(kb)" :precision="0" :value="server.network_tx / 1024" style="margin-right: 50px" />
                    </a-col>
                    <a-col :span="12">
                    <a-statistic title="下载(kb)" :precision="0" :value="server.network_rx / 1024" />
                    </a-col>
                </a-row>

            </div>


        </a-card>
      </div>
  
      <a-modal v-model:visible="isModalVisible" title="Server Details" @ok="handleOk" @cancel="handleCancel">
        <div v-if="selectedServer">
          <a-button type="primary" @click="deleteServer(selectedServer.name)" danger>Delete</a-button>
          <p><strong>Location:</strong> {{ selectedServer.location }} {{ selectedServer.name }}</p>
          <p><strong>Type:</strong> {{ selectedServer.type }}</p>
          <p><strong>Host:</strong> {{ selectedServer.host }}</p>
          <p><strong>Online (IPv4):</strong> {{ selectedServer.online4 ? 'Yes' : 'No' }}</p>
          <p><strong>Online (IPv6):</strong> {{ selectedServer.online6 ? 'Yes' : 'No' }}</p>
          <p><strong>Uptime:</strong> {{ formatUptime(selectedServer.uptime) }}</p>
          <!-- <p><strong>GPUInfo:</strong> {{ selectedServer.gpu_info }}</p> -->
          <div v-for="gpu in selectedServer.gpu_info">
            <a-row>
                    <!-- 要保证20 + 4 = 24. 根据这个来分配比例 -->
                    <a-col :span="18" :push="6">
                    
                        <a-progress
                            :stroke-color="{
                                '0%': '#108ee9',
                                '100%': '#87d068',
                            }"
                            :style="{width: '100%'}"
                            :percent="(gpu.GPUUtilization).toFixed(1)"
                            :format="percent => `${percent}%`"
                            />
                    </a-col>
                    <a-col :span="6" :pull="18">
                        GPU{{gpu.GPUIndex}} 占用：
                    </a-col>
              </a-row>
              <a-row>
                    <!-- 要保证20 + 4 = 24. 根据这个来分配比例 -->
                    <a-col :span="18" :push="6">
                    
                        <a-progress
                            :stroke-color="{
                                '0%': '#108ee9',
                                '100%': '#87d068',
                            }"
                            :style="{width: '100%'}"
                            :percent="(gpu.MemoryUtilization).toFixed(1)"
                            :format="percent => `${percent}%`"
                            />
                    </a-col>
                    <a-col :span="6" :pull="18">
                      GPU{{gpu.GPUIndex}} 内存：
                    </a-col>
              </a-row>
          </div>

          <p><strong>Load (1 min):</strong> {{ selectedServer.load_1 }}</p>
          <p><strong>Load (5 min):</strong> {{ selectedServer.load_5 }}</p>
          <p><strong>Load (15 min):</strong> {{ selectedServer.load_15 }}</p>
          <p><strong>Ping (10010):</strong> {{ selectedServer.ping_10010 }}</p>
          <p><strong>Ping (189):</strong> {{ selectedServer.ping_189 }}</p>
          <p><strong>Ping (10086):</strong> {{ selectedServer.ping_10086 }}</p>
          <p><strong>Time (10010):</strong> {{ selectedServer.time_10010 }} ms</p>
          <p><strong>Time (189):</strong> {{ selectedServer.time_189 }} ms</p>
          <p><strong>Time (10086):</strong> {{ selectedServer.time_10086 }} ms</p>
          <p><strong>TCP Connections:</strong> {{ selectedServer.tcp }}</p>
          <p><strong>UDP Connections:</strong> {{ selectedServer.udp }}</p>
          <p><strong>Processes:</strong> {{ selectedServer.process }}</p>
          <p><strong>Threads:</strong> {{ selectedServer.thread }}</p>
          <p><strong>Network RX:</strong> {{ (selectedServer.network_rx / 1024).toFixed(1) }} KB</p>
          <p><strong>Network TX:</strong> {{ (selectedServer.network_tx / 1024).toFixed(1) }} KB</p>
          <p><strong>Network In:</strong> {{ formatBytes(selectedServer.network_in) }}</p>
          <p><strong>Network Out:</strong> {{ formatBytes(selectedServer.network_out) }}</p>
          <p><strong>CPU Usage:</strong> {{ selectedServer.cpu }}%</p>
          <p><strong>Memory Total:</strong> {{ formatBytes(selectedServer.memory_total) }}</p>
          <p><strong>Memory Used:</strong> {{ formatBytes(selectedServer.memory_used) }}</p>
          <p><strong>Swap Total:</strong> {{ formatBytes(selectedServer.swap_total) }}</p>
          <p><strong>Swap Used:</strong> {{ formatBytes(selectedServer.swap_used) }}</p>
          <p><strong>HDD Total:</strong> {{ (selectedServer.hdd_total / 1024).toFixed(1) }} GB</p>
          <p><strong>HDD Used:</strong> {{ (selectedServer.hdd_used / 1024).toFixed(1)}} GB</p>
          <p><strong>All Network In:</strong> {{ formatBytes(selectedServer.network_in) }}</p>
          <p><strong>All Network Out:</strong> {{ formatBytes(selectedServer.network_out) }}</p>
          <p><strong>IO Read:</strong> {{ formatBytes(selectedServer.io_read)  }} TB</p>
          <p><strong>IO Write:</strong> {{ formatBytes(selectedServer.io_write)  }} TB</p>
          <p v-html="selectedServer.custom"></p>
        </div>
      </a-modal>
    </div>
  </template>
  <script>
  import axios from 'axios';
  import { StarOutlined, StarFilled, StarTwoTone } from '@ant-design/icons-vue';
  import { mapGetters } from 'vuex'


  export default {
    data() {
      return {
        isModalVisible: false,
        selectedServer: null,
        servers: []
        // backendServerUrl: "/"
      };
    },

    computed: {
    ...mapGetters(['backendServerUrl']),
    backendServerUrl() {
      return this.$store.state.backendServerUrl;
    }
  },
    created() {
          console.log("前置处理函数")
          console.log(this.backendServerUrl);
    },
    mounted() {
      this.fetchData();
      this.interval = setInterval(this.fetchData, 3000);
    },
    beforeDestroy() {
      clearInterval(this.interval);
    },
    methods: {


      formatUptime(seconds) {
      const days = Math.floor(seconds / (24 * 3600));
      let remainingSeconds = seconds % (24 * 3600);
      const hours = Math.floor(remainingSeconds / 3600);
      remainingSeconds %= 3600;
      const minutes = Math.floor(remainingSeconds / 60);
      const secs = remainingSeconds % 60;

      return `${days}d ${hours}h ${minutes}m ${secs}s`;
    },

      fetchData() {
        axios.get(  this.backendServerUrl + '/json/stats.json')
          .then(response => {
            this.servers = response.data.servers;
            // 按照servers的name排序
            this.servers.sort((a, b) => a.name.localeCompare(b.name));
          })
          .catch(error => {
            console.error('Error fetching data:', error);
          });
      },
      showModal(server) {
        this.selectedServer = server;
        this.isModalVisible = true;
        console.log(this.isModalVisible);
      },
      handleOk() {
        this.isModalVisible = false;
      },
      handleCancel() {
        this.isModalVisible = false;
      },
      deleteServer(serverName) {
        this.isModalVisible = false;
        // alert(`Server ${serverName} deleted`);
        // delete server here
        axios.get( this.backendServerUrl + '/server/del?target=' + serverName)
          .then(response => {
            alert(response.data + '')
          })
          .catch(error => {
            console.error('Error fetching data:', error);
            alert('Error delete server:', error);
          });

      },
      formatBytes(bytes) {
        if (bytes === 0) return '0 B';
        const k = 1024;
        const sizes = ['B', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
      }
    }
  };
  </script>

  
  <style>
  .custom-progress .ant-progress-text {
    font-size: 5px; /* 调整字体大小，示例为16px，可根据需要修改 */
  }

  a-space {
    display: flex;
    flex-direction: column;
    align-items: center;
    margin: 10px;
  }
  </style>