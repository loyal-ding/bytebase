<template>
  <span class="textlabel relative inline-flex items-center">
    <template v-if="!editable">
      {{ value || $t("label.empty-label-value") }}
    </template>
    <template v-else>
      <select
        class="absolute inset-0 opacity-0 m-0 p-0"
        :value="state.value"
        @change="onChange"
      >
        <option :value="''" :selected="!!state.value">
          {{ $t("label.empty-label-value") }}
        </option>
        <option
          v-for="(labelValue, i) in label.valueList"
          :key="i"
          :value="labelValue"
        >
          {{ labelValue }}
        </option>
      </select>
      {{ state.value || $t("label.empty-label-value") }}
      <heroicons-outline:chevron-down class="w-4 h-4 ml-0.5" />
    </template>
  </span>
</template>

<script lang="ts" setup>
import { capitalize } from "lodash-es";
import { useDialog } from "naive-ui";
import { computed, defineProps, defineEmits, reactive, watch } from "vue";
import { useI18n } from "vue-i18n";
import { Database, Label, LabelValueType } from "../../types";
import { hidePrefix } from "../../utils";

const props = defineProps<{
  label: Label;
  labelList: Label[];
  requiredLabelList: Label[];
  value: LabelValueType | undefined;
  database: Database;
  allowEdit: boolean;
}>();

const emit = defineEmits<{
  (e: "update:value", value: LabelValueType): void;
}>();

const { t } = useI18n();

const state = reactive({
  value: props.value,
});

const dialog = useDialog();

watch(
  () => props.value,
  (value) => (state.value = value)
);

const editable = computed((): boolean => {
  if (!props.allowEdit) return false;

  // Not editable if this is a required label in `dbNameTemplate`
  // e.g. tenant in "{{DB_NAME}}_{{TENANT}}"
  const isRequired = !!props.requiredLabelList.find(
    (label) => label.key === props.label.key
  );
  return !isRequired;
});

const onChange = (e: Event) => {
  const select = e.target as HTMLSelectElement;
  dialog.create({
    positiveText: t("common.confirm"),
    negativeText: t("common.cancel"),
    title: t("label.confirm-change", {
      label: capitalize(hidePrefix(props.label.key)),
    }),
    closable: false,
    maskClosable: false,
    closeOnEsc: false,
    onNegativeClick: () => {
      state.value = props.value;
      select.value = state.value || "";
    },
    onPositiveClick: () => {
      const { value } = select;
      emit("update:value", value);
    },
  });
};
</script>
