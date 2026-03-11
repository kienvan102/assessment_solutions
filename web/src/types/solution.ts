export interface InputField {
  name: string;
  label: string;
  type: 'string' | 'number' | 'boolean';
  required: boolean;
  default?: string;
  placeholder?: string;
}

export interface Solution {
  id: string;
  title: string;
  description: string;
  expected_behaviors: string[];
  code: string;
  endpoint?: string;
  method?: string;
  input_fields?: InputField[];
  sample_payload?: string;
}

export type RunRecord = {
  at: string;
  request: {
    endpoint: string;
    method: string;
    body?: string;
  };
  response: {
    status: number;
    data: any;
  };
};
