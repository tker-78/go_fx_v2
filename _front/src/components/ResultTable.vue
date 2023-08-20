<!-- apiからresultsの情報を取ってきて表示 -->

<template>
  <div class="v-container">
  <label>Results</label>
  <v-table class="result-table">
    <thead>
      <tr>
        <th class="text-center">Entry</th>
        <th class="text-center">Exit</th>
        <th class="text-center">差益</th>
        <th class="text-center">Swap利益</th>
        <th class="text-center">取引期間</th>
      </tr>
    </thead>

    <tbody>
      <tr v-for="r in results" :key="r.entry">
        <td>{{  r.entry }}</td>
        <td>{{  r.exit }}</td>
        <td>{{  r.capital_profit }}</td>
        <td>{{  r.swap_profit }}</td>
        <td>{{  r.duration }}</td>
      </tr>
    </tbody>
  </v-table>
  </div>

</template>

<script>
  export default{
    name: 'ResultTable',
    data() {
      return {
        result: {
          entry: '',
          exit: '',
          capital_profit: 0,
          swap_profit: 0,
          duration: ''
        },
      }
    },
    props: {
      results: {
        type: Object,
        default: () => {},
      }
    },
  }
</script>

<style>

</style>