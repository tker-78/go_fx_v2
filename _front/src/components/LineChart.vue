<template>
  <GChart
    type="LineChart"
    :data="chartData"
    :options="chartOptions"
  ></GChart>

</template>

<script>
  import { GChart } from 'vue-google-charts'

  export default {
    components: {
      GChart
    },
    props: {
      chartData: {
        type: Array,
        default: () => {
          return []
        }
      },
      chartOptions: {
        type: Object, 
        default: () => {
          return {}
        }
      }
    }
  }
</script>

<style scoped >
  template {
    width: 100px;
  }

</style>