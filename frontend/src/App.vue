<template>
  <div id="app">
    <DNSLookupTool/>
  </div>
</template>

<script>
import DNSLookupTool from './components/DNSLookupTool'

export default {
  name: 'App',
  components: {
    DNSLookupTool
  }
}
</script>

<style>
#app {
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  margin-top: 30px;
}
</style>
