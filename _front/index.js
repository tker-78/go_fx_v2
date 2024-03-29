const vuetify = Vuetify.createVuetify()


// Vue.component("gchart", VueGoogleCharts.GChart);



const app = Vue.createApp({
  data() {
    return {
      message: "Candle Chart",
      candles: [
      ['Mon', 20, 28, 38, 45],
      ['Tue', 31, 38, 55, 66],
      ['Wed', 50, 55, 77, 80],
      ['Thu', 77, 77, 66, 50],
      ['Fri', 68, 66, 22, 15]
      ]
    }
  },
});



app.use(vuetify).mount("#app")