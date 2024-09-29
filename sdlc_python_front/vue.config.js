module.exports = {
    configureWebpack: {
      resolve: {
        alias: {
          '@': require('path').resolve(__dirname, 'src')
        }
      }
    }
  };