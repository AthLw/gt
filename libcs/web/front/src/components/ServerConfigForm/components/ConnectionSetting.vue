<template>
  <el-form ref="ConnectionSettingRef" :model="localSetting" :rules="rules">
    <div class="card content-box">
      <el-descriptions :column="2" :border="true">
        <template #title> {{ $t("sconfig.ConnectionSetting") }} </template>
        <el-descriptions-item>
          <template #label>
            {{ $t("sconfig.Speed") }}
            <UsageTooltip :usage-text="$t('susage[\'Speed\']')" />
          </template>
          <el-form-item prop="Speed">
            <el-input-number v-model="localSetting.Speed" />
          </el-form-item>
        </el-descriptions-item>
        <el-descriptions-item>
          <template #label>
            {{ $t("sconfig.Connections") }}
            <UsageTooltip :usage-text="$t('susage[\'Connections\']')" />
          </template>
          <el-form-item prop="Connections">
            <el-input-number v-model="localSetting.Connections" />
          </el-form-item>
        </el-descriptions-item>
        <el-descriptions-item>
          <template #label>
            {{ $t("sconfig.ReconnectTimes") }}
            <UsageTooltip :usage-text="$t('susage[\'ReconnectTimes\']')" />
          </template>
          <el-form-item prop="ReconnectTimes">
            <el-input-number v-model="localSetting.ReconnectTimes" />
          </el-form-item>
        </el-descriptions-item>
        <el-descriptions-item>
          <template #label>
            {{ $t("sconfig.ReconnectDuration") }}
            <UsageTooltip :usage-text="$t('susage[\'ReconnectDuration\']')" />
          </template>
          <el-form-item prop="ReconnectDuration">
            <el-input v-model="localSetting.ReconnectDuration" />
          </el-form-item>
        </el-descriptions-item>
        <el-descriptions-item>
          <template #label>
            {{ $t("sconfig.Timeout") }}
            <UsageTooltip :usage-text="$t('susage[\'Timeout\']')" />
          </template>
          <el-form-item prop="Timeout">
            <el-input v-model="localSetting.Timeout" />
          </el-form-item>
        </el-descriptions-item>
        <el-descriptions-item>
          <template #label>
            {{ $t("sconfig.TimeoutOnUnidirectionalTraffic") }}
            <UsageTooltip :usage-text="$t('susage[\'TimeoutOnUnidirectionalTraffic\']')" />
          </template>
          <el-form-item prop="TimeoutOnUnidirectionalTraffic">
            <el-switch v-model="localSetting.TimeoutOnUnidirectionalTraffic" />
          </el-form-item>
        </el-descriptions-item>
      </el-descriptions>
    </div>
  </el-form>
</template>
<script setup name="ConnectionSetting" lang="ts">
import { FormInstance, FormRules } from "element-plus";
import { reactive, ref, watchEffect } from "vue";
import { ServerConfig } from "../interface";
import UsageTooltip from "@/components/UsageTooltip/index.vue";
import { validatorTimeFormat } from "@/utils/eleValidate";
import i18n from "@/languages";

interface ConnectionSettingProps {
  setting: ServerConfig.ConnectionSetting;
}

const props = withDefaults(defineProps<ConnectionSettingProps>(), {
  setting: () => ServerConfig.defaultConnectionSetting
});
const localSetting = reactive<ServerConfig.ConnectionSetting>({ ...props.setting });

//Sync with parent: props.setting -> localSetting
watchEffect(() => {
  Object.assign(localSetting, props.setting);
});

const emit = defineEmits(["update:setting"]);
//Sync with parent: localSetting -> emit("update:setting")
watchEffect(() => {
  emit("update:setting", localSetting);
});

//Form Related
const ConnectionSettingRef = ref<FormInstance>();
const rules = reactive<FormRules<ServerConfig.ConnectionSetting>>({
  Speed: [{ type: "number", message: "Please enter a number", trigger: "blur" }],
  Connections: [{ type: "number", message: "Please enter a number", trigger: "blur" }],
  ReconnectTimes: [{ type: "number", message: "Please enter a number", trigger: "blur" }],
  ReconnectDuration: [{ validator: validatorTimeFormat, trigger: "blur" }],
  Timeout: [{ validator: validatorTimeFormat, trigger: "blur" }]
});

const validateForm = (): Promise<void> => {
  return new Promise((resolve, reject) => {
    if (ConnectionSettingRef.value) {
      ConnectionSettingRef.value.validate(valid => {
        if (valid) {
          resolve();
        } else {
          reject(new Error(i18n.global.t("serror.ConnectionSettingValidationFailed")));
        }
      });
    } else {
      reject(new Error(i18n.global.t("serror.ConnectionSettingNotReady")));
    }
  });
};
defineExpose({
  validateForm
});
</script>
<style scoped lang="scss">
@import "../index.scss";
</style>
