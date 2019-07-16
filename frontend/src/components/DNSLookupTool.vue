<template>
  <div class="container notranslate">
        <div class="header clearfix">
            <nav>
                <ul class="nav nav-pill float-right">
                    <li class="nav-item">
                        Websocket Status:
                        <button @click="disconnect" v-if="status === 'connected'" class="btn btn-success btn-sm">Connected</button>
                        <button @click="connect" v-if="status === 'disconnected'" class="btn btn-secondary btn-sm">Disconnected</button>
                    </li>
                </ul>
            </nav>
            <div align="left">
              <div style="display: inline-block"><h2 class="text-muted"><a href="/"></a>DNS Lookup</h2></div>
              <div style="display: inline-block" class="spinner-grow spinner-grow-sm" role="status" v-show="typing"></div>
            </div>
        </div>

        <div class="jumbotron">
              <h5>Enter a domain name to view it's DNS records</h5>
              <small><b>Examples:</b>
                <a href="#" @click.prevent="message ='netflix.com'"> netflix.com, </a>
                <a href="#" @click.prevent="message ='google.com'">google.com, </a>
                <a href="#" @click.prevent="message ='reddit.com'">reddit.com, </a>
                <a href="#" @click.prevent="message ='yahoo.com'">yahoo.com</a><br>
              </small>

              <div v-if="status === 'connected'">
                <br>
                <form @submit.prevent="sendMessage" action="/">
                  <input v-model="message" class="form-control" placeholder="DNS Name..." @input="typing = true">

                </form>
                    <table class="table table-striped table-bordered">
                    <thead>
                        <tr>
                            <th>
                                Type
                            </th>
                            <th>
                                Domain Name
                            </th>
                            <th>
                                Address
                            </th>
                        </tr>
                    </thead>
                    <tbody>
                       <tr v-for="log in this.logs.A" v-bind:key="log">
                        <td >
                            <a>A</a>
                        </td>
                        <td >
                            <a>{{ logs.dnsname }}</a>
                        </td>
                        <td >
                            <a v-text="log"></a>
                        </td>
                        </tr>

                       <tr v-for="log in logs.AAAA" v-bind:key="log">
                        <td >
                            <a>AAAA</a>
                        </td>
                        <td >
                            <a>{{ logs.dnsname }}</a>
                        </td>
                        <td >
                            <a v-text="log"></a>
                        </td>
                        </tr>

                       <tr v-for="log in logs.MX" v-bind:key="log">
                        <td >
                            <a>MX</a>
                        </td>
                        <td >
                            <a>{{ logs.dnsname }}</a>
                        </td>
                        <td >
                            <a v-text="log"></a>
                        </td>
                        </tr>

                       <tr v-for="log in logs.NS" v-bind:key="log">
                        <td >
                            <a>NS</a>
                        </td>
                        <td >
                            <a>{{ logs.dnsname }}</a>
                        </td>
                        <td >
                            <a v-text="log"></a>
                        </td>
                        </tr>
                    </tbody>
                </table>
              </div>
        </div>
  </div>
</template>

<script>
import _ from 'lodash'
export default {
  name: 'DNSLookupTool',
  data () {
    return {
      message: '',
      typing: false,
      logs: {},
      status: 'disconnected',
      websocket_url: process.env.WEBSOCKET_URL
    }
  },

  watch: {
    message: _.debounce(function () {
      this.typing = false
      if (!this.message) {
        this.logs = {}
        return
      }
      this.sendMessage(this.message)
    }, 1600)
  },

  methods: {
    connect () {
      console.log(this.websocket_url)
      this.socket = new WebSocket(this.websocket_url)
      this.socket.onerror = error => {
        console.log(`WebSocket error: ${error}`)
        this.disconnect()
      }
      this.socket.onopen = () => {
        this.status = 'connected'
        console.log('WebSocket connected to:', this.socket.url)
        this.socket.onmessage = ({data}) => {
          this.logs = JSON.parse(data)
        }
      }
    },

    disconnect () {
      this.socket.close()
      this.status = 'disconnected'
      this.logs = ''
    },

    sendMessage (e) {
      this.socket.send(JSON.stringify({dnsname: this.message}))
      console.log('DNS query via Websocket:', this.message)
    }
  },

  created () {
    this.connect()
  }
}
</script>

<style scoped>
@import url('https://fonts.googleapis.com/css?family=Roboto|VT323');
h1 {
    font-family: 'Proza Libre', sans-serif;
    color: seagreen;
    font-weight: 300;
}

body {
  padding-top: 1.5rem;
  padding-bottom: 1.5rem;
  font-family: 'Roboto', sans-serif;
}

a {
    color: rgb(27, 27, 27);
}

.status {
    color:  white;
}

/* Custom page header */
.header {
  padding-bottom: 1rem;
  border-bottom: .05rem solid #e5e5e5;
  font-family: 'VT323', monospace;
}
/* Make the masthead heading the same height as the navigation */
.header h3 {
  margin-top: 0;
  margin-bottom: 0;
  line-height: 3rem;
}

/* Customize container */
@media (min-width: 48em) {
  .container {
    max-width: 46rem;
  }
}
.container-narrow > hr {
  margin: 2rem 0;
}

.footer {
    position: fixed;
    left: 0;
    bottom: 0;
    height: 40px;
    width: 100%;
    background-color: black;
    color: white;
    text-align: left;
}

</style>
