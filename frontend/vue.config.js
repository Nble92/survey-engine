module.exports = {
  chainWebpack: config => {
    if (process.env.NODE_ENV === 'development') {
      config
        .output
        .filename('[name].[hash].js') 
        .end() 
    }  
  },
  runtimeCompiler: true,
  
  devServer: {
      host: '0.0.0.0',
      port: 8080,
      disableHostCheck: true,
      proxy: {
        '/config': {target: 'http://127.0.0.1:8088'},
        '/v1/*': {target: 'http://127.0.0.1:8088'},
        '/functions.js': {target: 'http://127.0.0.1:8088'},
        '/instructionHTML/*': {target: 'http://127.0.0.1:8088'},
        '/images/*': {target: 'http://127.0.0.1:8088'},
      }
    }
}