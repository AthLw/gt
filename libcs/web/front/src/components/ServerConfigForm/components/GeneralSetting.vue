<template>
  <!-- General Setting -->
  <el-form ref="generalSettingRef" :model="localSetting" :rules="rules">
    <div class="card content-box">
      <el-descriptions :column="2" :border="true">
        <template #title>{{ $t("sconfig.GeneralSetting") }}</template>
        <!-- Users -->
        <el-descriptions-item>
          <template #label>
            {{ $t("sconfig.Users") }}
            <UsageTooltip :usage-text="$t('susage[\'Users\']')" />
          </template>
          <el-form-item prop="Users">
            <el-input v-model="localSetting.UserPath"></el-input>
          </el-form-item>
        </el-descriptions-item>
        <!-- AuthAPI -->
        <el-descriptions-item>
          <template #label>
            {{ $t("sconfig.AuthAPI") }}
            <UsageTooltip :usage-text="$t('susage[\'AuthAPI\']')" />
          </template>
          <el-form-item prop="AuthAPI">
            <el-input v-model="localSetting.AuthAPI"></el-input>
          </el-form-item>
        </el-descriptions-item>
        <!-- TCP Number -->
        <el-descriptions-item>
          <template #label>
            {{ $t("sconfig.TCPNumber") }}
            <UsageTooltip :usage-text="$t('susage[\'TCPNumber\']')" />
          </template>
          <el-form-item prop="TCPNumber">
            <el-input-number v-model="localSetting.TCPNumber" :min="0" />
          </el-form-item>
        </el-descriptions-item>
        <!-- WebAddr -->
        <el-descriptions-item>
          <template #label>
            {{ $t("sconfig.WebAddr") }}
            <UsageTooltip :usage-text="$t('susage[\'WebAddr\']')" />
          </template>
          <el-form-item prop="WebAddr">
            <el-input v-model="localSetting.WebAddr"></el-input>
          </el-form-item>
        </el-descriptions-item>
      </el-descriptions>
      <el-row style="width: 100%">
        <el-collapse style="width: 100%">
          <el-collapse-item name="1">
            <template #title>
              <el-text style="width: 100%" size="large">{{ $t("cconfig.DetailSettings") }} </el-text>
            </template>
            <el-row :gutter="10" style="width: 100%">
              <el-col :span="12">
                <TCPSetting ref="tcpSettingRef" :setting="tcpSetting" @update:setting="updateTCPSetting" />
              </el-col>
              <el-col :span="12">
                <HostSetting ref="hostSettingRef" :setting="hostSetting" @update:setting="updateHostSetting" />
              </el-col>
            </el-row>
          </el-collapse-item>
        </el-collapse>
      </el-row>
    </div>
  </el-form>
</template>

<script setup lang="ts" name="GeneralSetting">
import UsageTooltip from "@/components/UsageTooltip/index.vue";
import { ServerConfig } from "../interface";
import { reactive, ref, watch } from "vue";
import { FormInstance, FormRules } from "element-plus";
import TCPSetting from "./TCPSetting.vue";
import HostSetting from "./HostSetting.vue";
import cloneDeep from "lodash/cloneDeep";
import i18n from "@/languages";

interface GeneralSettingProps {
  setting: ServerConfig.GeneralSettingProps;
}

const props = withDefaults(defineProps<GeneralSettingProps>(), {
  setting: () => ServerConfig.getDefaultGeneralSettingProps()
});

//use deep clone to avoid changing props
const localSetting = reactive<ServerConfig.GeneralSettingProps>(cloneDeep(props.setting));

//use shallow clone to avoid sync in the current component
const tcpSetting = reactive<ServerConfig.TCP[]>(localSetting.TCPs);
const hostSetting = reactive<ServerConfig.Host>(localSetting.Host);

//Sync with parent: props.setting -> localSetting(tcpSetting, hostSetting)
watch(
  () => props.setting,
  () => {
    localSetting.UserPath = props.setting.UserPath;
    localSetting.AuthAPI = props.setting.AuthAPI;
    localSetting.TCPNumber = props.setting.TCPNumber;
    localSetting.WebAddr = props.setting.WebAddr;
    tcpSetting.splice(0, tcpSetting.length, ...props.setting.TCPs);
    hostSetting.Number = props.setting.Host.Number;
    hostSetting.RegexStr.splice(0, hostSetting.RegexStr.length, ...props.setting.Host?.RegexStr);
    hostSetting.WithID = props.setting.Host.WithID;
  },
  { deep: true }
);

const emit = defineEmits(["update:setting"]);

//Sync with parent: localSetting(tcpSetting, hostSetting) -> emit("update:setting")
watch(
  () => localSetting,
  () => {
    emit("update:setting", localSetting);
  },
  { deep: true }
);

//Form Related
const generalSettingRef = ref<FormInstance>();
const tcpSettingRef = ref<InstanceType<typeof TCPSetting> | null>(null);
const hostSettingRef = ref<InstanceType<typeof HostSetting> | null>(null);

const rules = reactive<FormRules>({});

//Sync with child
const updateTCPSetting = (setting: ServerConfig.TCP[]) => {
  tcpSetting.splice(0, tcpSetting.length, ...setting);
};
const updateHostSetting = (setting: ServerConfig.Host) => {
  if (JSON.stringify(hostSetting) === JSON.stringify(setting)) return;
  hostSetting.Number = setting.Number;
  hostSetting.RegexStr.splice(0, hostSetting.RegexStr.length, ...setting.RegexStr);
  hostSetting.WithID = setting.WithID;
};

const validateForm = (): Promise<void> => {
  const validations = [
    tcpSettingRef.value?.validateForm(),
    hostSettingRef.value?.validateForm(),
    new Promise<void>((resolve, reject) => {
      if (generalSettingRef.value) {
        generalSettingRef.value.validate(valid => {
          if (valid) {
            resolve();
          } else {
            reject(new Error(i18n.global.t("serror.GeneralSettingValidationFailed")));
          }
        });
      } else {
        reject(new Error(i18n.global.t("serror.GeneralSettingNotReady")));
      }
    })
  ];
  return Promise.all(validations).then(() => {
    console.log(i18n.global.t("serror.GeneralSettingValidationPassed"));
  });
};

defineExpose({
  validateForm
});
</script>
<style scoped lang="scss">
@import "../index.scss";
.el-row {
  margin-bottom: 20px;
}
.el-row:last-child {
  margin-bottom: 0;
}
.el-col {
  border-radius: 4px;
}
</style>
