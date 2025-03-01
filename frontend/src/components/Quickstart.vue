<template>
  <div class="space-y-2 w-full pl-4 pr-2">
    <div class="flex flex-row justify-between">
      <div class="outline-title toplevel flex">
        {{ $t("common.quickstart") }}
      </div>
      <button class="btn-icon" @click.prevent="hideQuickstart">
        <heroicons-solid:x class="w-4 h-4" />
      </button>
    </div>
    <nav class="flex justify-center" aria-label="Progress">
      <ol class="space-y-4 w-full">
        <li v-for="(intro, index) in effectiveList" :key="index">
          <!-- Complete Task -->
          <!-- use <router-link> if intro.link is not empty -->
          <!-- use <span> otherwise -->
          <component
            :is="intro.link ? 'router-link' : 'span'"
            :to="intro.link"
            class="group cursor-pointer"
            @click="intro.click"
          >
            <span class="flex items-start">
              <span
                class="flex-shrink-0 relative h-5 w-5 flex items-center justify-center"
              >
                <template v-if="intro.done">
                  <!-- Heroicon name: solid/check-circle -->
                  <heroicons-solid:check-circle
                    class="h-full w-full text-success group-hover:text-success-hover"
                  />
                </template>
                <template v-else-if="isTaskActive(index)">
                  <span
                    class="absolute h-4 w-4 rounded-full bg-blue-200"
                  ></span>
                  <span
                    class="relative block w-2 h-2 bg-info rounded-full"
                  ></span>
                </template>
                <template v-else>
                  <div
                    class="h-2 w-2 bg-gray-300 rounded-full group-hover:bg-gray-400"
                  ></div>
                </template>
              </span>
              <span
                class="ml-2 text-sm font-medium text-control-light group-hover:text-control-light-hover"
                :class="intro.done ? 'line-through' : ''"
                >{{ intro.name }}</span
              >
            </span>
          </component>
        </li>
      </ol>
    </nav>
  </div>
</template>

<script lang="ts">
import { computed, reactive, ComputedRef, defineComponent } from "vue";
import { useStore } from "vuex";
import { useI18n } from "vue-i18n";
import { isDBA, isOwner } from "../utils";
import { useKBarHandler, useKBarEventOnce } from "@bytebase/vue-kbar";

type IntroItem = {
  name: string;
  link: string;
  allowDBA: boolean;
  allowDeveloper: boolean;
  done: ComputedRef<boolean>;
  click?: () => void;
};

export default defineComponent({
  name: "QuickStart",
  setup() {
    const store = useStore();
    const { t } = useI18n();
    const kbarHandler = useKBarHandler();

    const currentUser = computed(() => store.getters["auth/currentUser"]());

    const introList: IntroItem[] = reactive([
      {
        name: computed(() => t("quick-start.bookmark-an-issue")),
        link: "/issue/hello-world-101",
        allowDBA: true,
        allowDeveloper: true,
        done: computed(() => {
          return store.getters["uistate/introStateByKey"]("bookmark.create");
        }),
      },
      {
        name: computed(() => t("quick-start.add-a-comment")),
        link: "/issue/hello-world-101#activity101",
        allowDBA: true,
        allowDeveloper: true,
        done: computed(() => {
          return store.getters["uistate/introStateByKey"]("comment.create");
        }),
      },
      {
        name: computed(() => t("quick-start.visit-project")),
        link: "/project",
        allowDBA: true,
        allowDeveloper: true,
        done: computed(() => {
          return store.getters["uistate/introStateByKey"]("project.visit");
        }),
      },
      {
        name: computed(() => t("quick-start.visit-environment")),
        link: "/environment",
        allowDBA: true,
        allowDeveloper: false,
        done: computed(() => {
          return store.getters["uistate/introStateByKey"]("environment.visit");
        }),
      },
      {
        name: computed(() => t("quick-start.visit-instance")),
        link: "/instance",
        allowDBA: true,
        allowDeveloper: !store.getters["subscription/feature"](
          "bb.feature.dba-workflow"
        ),
        done: computed(() => {
          return store.getters["uistate/introStateByKey"]("instance.visit");
        }),
      },
      {
        name: computed(() => t("quick-start.visit-database")),
        link: "/db",
        allowDBA: true,
        allowDeveloper: true,
        done: computed(() =>
          store.getters["uistate/introStateByKey"]("database.visit")
        ),
      },
      {
        name: computed(() => t("quick-start.add-a-member")),
        link: "/setting/member",
        allowDBA: false,
        allowDeveloper: false,
        done: computed(() =>
          store.getters["uistate/introStateByKey"]("member.addOrInvite")
        ),
      },
      {
        name: computed(() =>
          t("quick-start.use-kbar", {
            shortcut: `${navigator.platform.match(/mac/i) ? "cmd" : "ctrl"}-k`,
          })
        ),
        link: "",
        allowDBA: true,
        allowDeveloper: true,
        click: () => {
          kbarHandler.value.show();
        },
        done: computed(() =>
          store.getters["uistate/introStateByKey"]("kbar.open")
        ),
      },
    ]);

    const effectiveList = computed(() => {
      if (isOwner(currentUser.value.role)) {
        return introList;
      }
      if (isDBA(currentUser.value.role)) {
        return introList.filter((item) => item.allowDBA);
      }
      return introList.filter((item) => item.allowDeveloper);
    });

    const isTaskActive = (index: number): boolean => {
      for (let i = index - 1; i >= 0; i--) {
        if (!effectiveList.value[i].done) {
          return false;
        }
      }
      return !effectiveList.value[index].done;
    };

    const hideQuickstart = () => {
      store
        .dispatch("uistate/saveIntroStateByKey", {
          key: "hidden",
          newState: true,
        })
        .then(() => {
          store.dispatch("notification/pushNotification", {
            module: "bytebase",
            style: "INFO",
            title: computed(() => t("quick-start.notice.title")),
            description: computed(() => t("quick-start.notice.desc")),
            manualHide: true,
          });
        });
    };

    useKBarEventOnce("open", () => {
      // once kbar is open, mark the quickstart as done
      store.dispatch("uistate/saveIntroStateByKey", {
        key: "kbar.open",
        newState: true,
      });
    });

    return {
      effectiveList,
      isTaskActive,
      hideQuickstart,
    };
  },
});
</script>
