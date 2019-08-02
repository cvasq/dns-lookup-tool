<template>
  <div id="app">
    <DNSLookupTool/>
    <Footer/>
  </div>
</template>

<script>
import DNSLookupTool from './components/DNSLookupTool'
import Footer from './components/Footer'                                                                                                                                                      


export default {
  name: 'App',
  components: {
    DNSLookupTool,
    Footer
  }
}
</script>

<style>
#app {                                                                                                                                                                                        
  font-family: 'Avenir', Helvetica, Arial, sans-serif;                                                                                                                                        
  -webkit-font-smoothing: antialiased;                                                                                                                                                        
  -moz-osx-font-smoothing: grayscale;                                                                                                                                                         
  text-align: center;                                                                                                                                                                         
  color: #2c3e50;                                                                                                                                                                             
  margin-top: 60px;                                                                                                                                                                           
  padding-top: 1.5rem;                                                                                                                                                                        
  padding-bottom: 1.5rem;                                                                                                                                                                     
}  
</style>
