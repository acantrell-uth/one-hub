const defaultConfig = {
  input: {
    name: '',
    type: 1,
    key: '',
    base_url: '',
    other: '',
    proxy: '',
    test_model: '',
    model_mapping: [],
    models: [],
    groups: ['default'],
    plugin: {},
    tag: '',
    only_chat: false,
    pre_cost: 1
  },
  inputLabel: {
    name: '渠道名称',
    type: '渠道类型',
    base_url: '渠道API地址',
    key: '密钥',
    other: '其他参数',
    proxy: '代理地址',
    test_model: '测速模型',
    models: '模型',
    model_mapping: '模型映射关系',
    groups: '用户组',
    only_chat: '仅支持聊天',
    tag: '标签',
    provider_models_list: '',
    pre_cost: '预计费选项'
  },
  prompt: {
    type: '请选择渠道类型',
    name: '请为渠道命名',
    base_url: '可空，请输入中转API地址，例如通过cloudflare中转',
    key: '请输入渠道对应的鉴权密钥',
    other: '',
    proxy: '单独设置代理地址，支持http和socks5，例如：http://127.0.0.1:1080',
    test_model: '用于测试使用的模型，为空时无法测速,如：gpt-3.5-turbo，仅支持chat模型',
    models:
      '请选择该渠道所支持的模型,你也可以输入通配符*来匹配模型，例如：gpt-3.5*，表示支持所有gpt-3.5开头的模型，*号只能在最后一位使用，前面必须有字符，例如：gpt-3.5*是正确的，*gpt-3.5是错误的',
    model_mapping: '模型映射关系：例如用户请求A模型，实际转发给渠道的模型为B。',
    model_headers: '自定义模型请求头，例如：{"key": "value"}',
    groups: '请选择该渠道所支持的用户组',
    only_chat: '如果选择了仅支持聊天，那么遇到有函数调用的请求会跳过该渠道',
    provider_models_list: '必须填写所有数据后才能获取模型列表',
    tag: '你可以为你的渠道打一个标签，打完标签后，可以通过标签进行批量管理渠道，注意：设置标签后某些设置只能通过渠道标签修改，无法在渠道列表中修改。',
    pre_cost:
      '这里选择预计费选项，用于预估费用，如果你觉得计算图片占用太多资源，可以选择关闭图片计费。但是请注意：有些渠道在stream下是不会返回tokens的，这会导致输入tokens计算错误。'
  },
  modelGroup: 'OpenAI'
};

const typeConfig = {
  1: {
    inputLabel: {
      provider_models_list: '从OpenAI获取模型列表'
    }
  },
  3: {
    inputLabel: {
      base_url: 'AZURE_OPENAI_ENDPOINT',
      other: '默认 API 版本'
    },
    prompt: {
      base_url: '请填写AZURE_OPENAI_ENDPOINT',
      other: '请输入默认API版本，例如：2024-05-01-preview'
    }
  },
  14: {
    input: {
      models: [
        'claude-instant-1.2',
        'claude-2.0',
        'claude-2.1',
        'claude-3-opus-20240229',
        'claude-3-sonnet-20240229',
        'claude-3-haiku-20240307'
      ],
      test_model: 'claude-3-haiku-20240307'
    },
    modelGroup: 'Anthropic'
  },
  24: {
    inputLabel: {
      other: '位置/区域'
    },
    input: {
      models: ['tts-1', 'tts-1-hd']
    },
    prompt: {
      test_model: '',
      base_url: '',
      other: '请输入你 Speech Studio 的位置/区域，例如：eastasia'
    }
  },
  25: {
    inputLabel: {
      other: '版本号',
      provider_models_list: '从Gemini获取模型列表'
    },
    input: {
      models: ['gemini-pro', 'gemini-pro-vision', 'gemini-1.0-pro', 'gemini-1.5-pro'],
      test_model: 'gemini-pro'
    },
    prompt: {
      other: '请输入版本号，例如：v1'
    },
    modelGroup: 'Google Gemini'
  },
  30: {
    input: {
      models: [
        'open-mistral-7b',
        'open-mixtral-8x7b',
        'mistral-small-latest',
        'mistral-medium-latest',
        'mistral-large-latest',
        'mistral-embed'
      ],
      test_model: 'open-mistral-7b'
    },
    inputLabel: {
      provider_models_list: '从Mistral获取模型列表'
    },
    modelGroup: 'Mistral'
  },
  39: {
    input: {
      models: ['phi3', 'llama3']
    },
    prompt: {
      base_url: '请输入你部署的Ollama地址，例如：http://127.0.0.1:11434，如果你使用了cloudflare Zero Trust，可以在下方插件填入授权信息',
      key: '请随意填写'
    }
  },
};

export { defaultConfig, typeConfig };
