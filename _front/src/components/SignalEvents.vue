<!-- apiからsignalEventsの情報を取ってきて表示 -->

<template>
  <div class="v-container">
    <label>Signals</label>
    <v-table class="signal-table ">
      <thead>
        <tr>
          <th class="text-center">time</th>
          <th class="text-center">currency_code</th>
          <th class="text-center">side</th>
          <th class="text-center">price</th>
          <th class="text-center">size</th>
        </tr>
      </thead>

      <tbody>
        <tr v-for="s in signals" :key="s.time">
          <td> {{ s.time }}</td> 
          <td>{{  s.currency_code }}</td>
          <td>{{  s.side }}</td>
          <td>{{  s.price }}</td>
          <td>{{  s.size }}</td>
        </tr>

      </tbody>
      <br>
      <br>


    </v-table>
  </div>

</template>

<script>
  export default {
    name: 'SignalEvents',
    data() {
      return {
        signal: {
          time: "", 
          currency_code: "",
          side: "",
          price: 0,
          size: 0
        },
      }
    },
    props: {
      signals: {
        type: Array,
        default: () => [],
      }
    },
  }
</script>

<style scoped>

</style>