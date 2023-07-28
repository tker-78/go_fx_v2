const vuetify = Vuetify.createVuetify()


const app = Vue.createApp({
  data() {
    return {
      message: "Candle Chart",
      candles: [],
    }
  },
});



app.use(vuetify).mount("#app")