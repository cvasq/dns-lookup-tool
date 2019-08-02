module.exports = {                                                                                                                                                                                                                                                                                                           
} 

var webpack = require('webpack')
	
module.exports = {
  publicPath: process.env.VUE_APP_BASE_PATH === undefined || process.env.VUE_APP_BASE_PATH === null                                                                                                                                                                                                                          
    ? '/'
    : process.env.VUE_APP_BASE_PATH,                                                                                                                                                                                                                                                                                         
  runtimeCompiler: true,
  configureWebpack: {
    plugins: [
      new webpack.DefinePlugin({
        'process.env': {
          'VUE_APP_WS_URL': process.env.VUE_APP_WS_URL
        }
      })
    ]
  }
}
