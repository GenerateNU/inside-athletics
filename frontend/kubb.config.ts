import { defineConfig } from '@kubb/core'
import { pluginOas } from '@kubb/plugin-oas'
import { pluginTs } from '@kubb/plugin-ts'
import { pluginClient } from '@kubb/plugin-client'
import { pluginReactQuery } from '@kubb/plugin-react-query'
import { pluginZod } from '@kubb/plugin-zod'

export default defineConfig({
  root: '.',
  input: {
    path: '../openapi.yaml',
  },
  output: {
    path: './api',
    clean: true,
  },
  plugins: [
    pluginOas(),
    pluginTs({
      output: {
        path: './models',
      },
    }),
    pluginClient({
      output: {
        path: './clients',
      },
    }),
    pluginReactQuery({
      output: {
        path: './hooks',
      },
    }),
    pluginZod({
      output: {
        path: './zod',
      },
    }),
  ],
})
