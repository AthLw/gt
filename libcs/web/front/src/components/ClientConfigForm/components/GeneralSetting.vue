<template>
  <!-- General Setting -->
  <el-form ref="generalSettingRef" :model="localSetting" :rules="rules">
    <div class="card content-box">
      <el-descriptions :column="2" :border="true">
        <template #title> {{ $t("cconfig.GeneralSetting") }} </template>
        <!-- ID -->
        <el-descriptions-item>
          <template #label>
            {{ $t("cconfig.ID") }}
            <UsageTooltip :usage-text="$t('cusage[\'ID\']')" />
          </template>
          <el-form-item prop="ID">
            <el-input v-model="localSetting.ID"></el-input>
          </el-form-item>
        </el-descriptions-item>
        <!-- Secret -->
        <el-descriptions-item>
          <template #label>
            {{ $t("cconfig.Secret") }}
            <UsageTooltip :usage-text="$t('cusage[\'Secret\']')" />
          </template>
          <el-form-item prop="Secret">
            <el-input v-model="localSetting.Secret" type="password" show-password></el-input>
          </el-form-item>
        </el-descriptions-item>
        <!-- ReconnectDelay -->
        <el-descriptions-item>
          <template #label>
            {{ $t("cconfig.ReconnectDelay") }}
            <UsageTooltip :usage-text="$t('cusage[\'ReconnectDelay\']')" />
          </template>
          <el-form-item prop="ReconnectDelay">
            <el-input v-model="localSetting.ReconnectDelay"></el-input>
          </el-form-item>
        </el-descriptions-item>
        <!-- RemoteTimeout -->
        <el-descriptions-item>
          <template #label>
            {{ $t("cconfig.RemoteTimeout") }}
            <UsageTooltip :usage-text="$t('cusage[\'RemoteTimeout\']')" />
          </template>
          <el-form-item prop="RemoteTimeout">
            <el-input v-model="localSetting.RemoteTimeout"></el-input>
          </el-form-item>
        </el-descriptions-item>
        <!-- Remote -->
        <template v-for="(remote, index) in localSetting.Remote" :key="index">
          <el-descriptions-item class="remote_des">
            <template #label>
              {{ $t("cconfig.Remote") }}
              <UsageTooltip :usage-text="$t('cusage[\'Remote\']')" />
            </template>
            <el-form-item prop="Remote">
              <el-select class="remote_select" placeholder="Protocol" v-model="localSetting.Remote[index].Protocol">
                <el-option v-for="protocol in protocols" :label="protocol.label" :value="protocol.value" :key="protocol.value" />
              </el-select>
              <el-input style="display: inline-block" v-model="localSetting.Remote[index].Addr" placeholder="127.0.0.1:8080" />
              <el-button @click="deleteRemote(index)">{{ $t("cconfig.Delete") }}</el-button>
              <el-button @click="addRemote()">{{ $t("cconfig.AddNewRemote") }}</el-button>
            </el-form-item>
          </el-descriptions-item>
        </template>
      </el-descriptions>
      <el-row style="width: 100%">
        <el-collapse style="width: 100%">
          <el-collapse-item>
            <template #title>
              <el-text size="large" style="width: 100%">{{ $t("cconfig.DetailSettings") }}</el-text>
            </template>
            <el-descriptions :column="2" :border="true">
              <!-- RemoteSTUN -->
              <template v-for="(stun, index) in localSetting.RemoteSTUN" :key="index">
                <el-descriptions-item>
                  <template #label>
                    {{ $t("cconfig.RemoteSTUN") }}
                    <UsageTooltip :usage-text="$t('cusage[\'RemoteSTUN\']')" />
                  </template>
                  <el-form-item prop="RemoteSTUN">
                    <el-input v-model="localSetting.RemoteSTUN[index]"></el-input>
                  </el-form-item>
                </el-descriptions-item>
              </template>
              <!-- RemoteAPI -->
              <el-descriptions-item>
                <template #label>
                  {{ $t("cconfig.RemoteAPI") }}
                  <UsageTooltip :usage-text="$t('cusage[\'RemoteAPI\']')" />
                </template>
                <el-form-item prop="RemoteAPI">
                  <el-input v-model="localSetting.RemoteAPI"></el-input>
                </el-form-item>
              </el-descriptions-item>
              <!-- WebAddr -->
              <el-descriptions-item>
                <template #label>
                  {{ $t("cconfig.WebAddr") }}
                  <UsageTooltip :usage-text="$t('cusage[\'WebAddr\']')" />
                </template>
                <el-form-item prop="WebAddr">
                  <el-input v-model="localSetting.WebAddr"></el-input>
                </el-form-item>
              </el-descriptions-item>
              <!-- RemoteCert -->
              <el-descriptions-item>
                <template #label>
                  {{ $t("cconfig.RemoteCert") }}
                  <UsageTooltip :usage-text="$t('cusage[\'RemoteCert\']')" />
                </template>
                <el-form-item prop="RemoteCert">
                  <el-input v-model="localSetting.RemoteCert"></el-input>
                </el-form-item>
              </el-descriptions-item>
              <!-- RemoteCertInsecure -->
              <el-descriptions-item>
                <template #label>
                  {{ $t("cconfig.RemoteCertInsecure") }}
                  <UsageTooltip :usage-text="$t('cusage[\'RemoteCertInsecure\']')" />
                </template>
                <el-switch
                  v-model="localSetting.RemoteCertInsecure"
                  :active-text="$t('cconfig.true')"
                  :inactive-text="$t('cconfig.false')"
                />
              </el-descriptions-item>
              <!-- RemoteConnections -->
              <el-descriptions-item>
                <template #label>
                  {{ $t("cconfig.RemoteConnections") }}
                  <UsageTooltip :usage-text="$t('cusage[\'RemoteConnections\']')" />
                </template>
                <el-form-item prop="RemoteConnections">
                  <el-input-number v-model="localSetting.RemoteConnections" :min="1" :max="10" />
                </el-form-item>
              </el-descriptions-item>
              <!-- RemoteIdleConnections -->
              <el-descriptions-item>
                <template #label>
                  {{ $t("cconfig.RemoteIdleConnections") }}
                  <UsageTooltip :usage-text="$t('cusage[\'RemoteIdleConnections\']')" />
                </template>
                <el-form-item prop="RemoteIdleConnections">
                  <el-input-number v-model="localSetting.RemoteIdleConnections" :min="0" :max="localSetting.RemoteConnections" />
                </el-form-item>
              </el-descriptions-item>
            </el-descriptions>
          </el-collapse-item>
        </el-collapse>
      </el-row>
    </div>
  </el-form>
</template>

<script setup name="GeneralSetting" lang="ts">
import { reactive, ref, watchEffect } from "vue";
import { ClientConfig } from "../interface";
import UsageTooltip from "@/components/UsageTooltip/index.vue";
import type { FormInstance, FormRules } from "element-plus";
import { validatorTimeFormat } from "@/utils/eleValidate";
import i18n from "@/languages";

interface GeneralSettingProps {
  setting: ClientConfig.GeneralSetting;
}

const props = withDefaults(defineProps<GeneralSettingProps>(), {
  setting: () => ClientConfig.defaultGeneralSetting
});
const localSetting = reactive<ClientConfig.GeneralSetting>({ ...props.setting });

const protocols = [
  {
    label: "TCP",
    value: "tcp"
  },
  {
    label: "TLS",
    value: "tls"
  }
];
//Sync with parent: props.setting -> localSetting
watchEffect(() => {
  Object.assign(localSetting, props.setting);
});

const emit = defineEmits(["update:setting"]);
//Sync with parent: localSetting -> emit("update:setting")
watchEffect(() => {
  emit("update:setting", localSetting);
});

const generalSettingRef = ref<FormInstance>();
const validatorRemoteIdleConnections = (rule: any, value: number, callback: any) => {
  if (value < 0 || value > localSetting.RemoteConnections) {
    callback(new Error(i18n.global.t("cerror.PleaseInputRemoteIdleConnections")));
  } else {
    callback();
  }
};
const validatorRemote = (index: number) => (rule: any, values: any, callback: any) => {
  console.log("Calling validatorRemote");
  const value = localSetting.Remote[index];
  if (!value.Protocol || value.Protocol.trim() === "" || !value.Addr || value.Addr.trim() === "") {
    callback(new Error(i18n.global.t("cerror.PleaseEnterValidRemote")));
  } else {
    callback();
  }
};
const validatorRemoteAPI = (rule: any, value: any, callback: any) => {
  console.log("Calling validatorRemoteAPI");
  if (!value) {
    callback();
  } else if (value.startsWith("http://") || value.startsWith("https://")) {
    console.log("Valid remoteAPI format");
    callback();
  } else {
    console.log("Invalid remoteAPI format");
    callback(new Error(i18n.global.t("cerror.PleaseEnterValidRemoteAPI")));
  }
};
const rules = reactive<FormRules<ClientConfig.GeneralSetting>>({
  ID: [{ required: true, message: i18n.global.t("cerror.PleaseInputID"), trigger: "blur" }],
  Secret: [{ message: i18n.global.t("cerror.PleaseInputSecret"), trigger: "blur" }],
  ReconnectDelay: [{ validator: validatorTimeFormat, trigger: "blur" }],
  RemoteTimeout: [{ validator: validatorTimeFormat, trigger: "blur" }],
  Remote: localSetting.Remote.map((_, index) => ({
    validator: validatorRemote(index),
    trigger: "change"
  })),
  RemoteAPI: [{ validator: validatorRemoteAPI, trigger: "blur" }],
  RemoteConnections: [
    { type: "number", message: i18n.global.t("cerror.PleaseInputRemoteConnections"), trigger: "blur" },
    {
      type: "number",
      min: 1,
      max: 10,
      message: i18n.global.t("cerror.PleaseInputRemoteConnectionsBetween1And10"),
      trigger: "blur"
    }
  ],
  RemoteIdleConnections: [
    { type: "number", message: i18n.global.t("cerror.PleaseInputRemoteIdleConnections"), trigger: "blur" },
    {
      validator: validatorRemoteIdleConnections,
      trigger: "change"
    }
  ]
});
const addRemote = () => {
  localSetting.Remote.push({ Protocol: "", Addr: "" });
};

const deleteRemote = (index: number) => {
  if (localSetting.Remote.length != 1) {
    localSetting.Remote.splice(index, 1);
  }
};

const checkRemoteSetting = (): Promise<void> => {
  return new Promise<void>((resolve, reject) => {
    const isRemoteEmpty = localSetting.Remote.every(item => !item.Addr.trim() && !item.Protocol.trim());
    const isRemoteAPIEmpty = !localSetting.RemoteAPI?.trim();

    if (isRemoteEmpty && isRemoteAPIEmpty) {
      reject(new Error(i18n.global.t("cerror.PleaseInputRemoteOrRemoteAPI")));
    } else {
      resolve();
    }
  });
};

const validateForm = (): Promise<void> => {
  const validations = [
    checkRemoteSetting(),
    new Promise<void>((resolve, reject) => {
      if (generalSettingRef.value) {
        generalSettingRef.value.validate(valid => {
          if (valid) {
            console.log(i18n.global.t("cerror.GeneralSettingValidationPassed"));
            resolve();
          } else {
            console.log(i18n.global.t("cerror.GeneralSettingValidationFailed"));
            reject(new Error(i18n.global.t("cerror.GeneralSettingValidationFailedCheckInput")));
          }
        });
      } else {
        reject(new Error(i18n.global.t("cerror.GeneralSettingNotReady")));
      }
    })
  ];
  return Promise.all(validations).then(
    () => {
      console.log(i18n.global.t("cerror.GeneralSettingValidationPassed"));
      return Promise.resolve();
    },
    error => {
      return Promise.reject(error);
    }
  );
};
defineExpose({
  validateForm
});
</script>

<style scoped lang="scss">
@import "../index.scss";
</style>
