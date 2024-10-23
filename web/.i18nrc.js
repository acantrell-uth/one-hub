const { defineConfig } = require('@lobehub/i18n-cli');

module.exports = defineConfig({
  entry: 'src/i18n/locales/en_US.json',
  entryLocale: 'en_US',
  output: 'src/i18n/locales',
  outputLocales: ['en_US', 'ja_JP'],
  temperature: 0,
  modelName: 'gpt-3.5-turbo',
  splitToken: 2048,
  experimental: {
    jsonMode: true
  }
});
