<template>
  <div>
    <div class="layout header">
      <h2>{{this.$t("invitation.mine")}}</h2>
      <p class="subtitle">
        <span>{{this.$t("invitation.used_friends")}}{{this.usedInvitations.length}}</span>
        <span
          class="pending"
        >{{this.$t("invitation.pending_friends")}}{{this.pendingInvitations.length}}</span>
      </p>
    </div>
    <van-tabs v-model="activeName">
      <van-tab :title="titleCodes" name="codes">
        <van-list v-model="loading" :finished="finished" :finished-text="finishedText">
          <InvitationCodeListItem
            class="layout"
            v-for="item in availableInvitations"
            :key="item.code"
            :invitation="item"
          />
        </van-list>
        <div class="layout">
          <p class="subtitle">{{this.$t("invitation.rule1")}}</p>
          <van-button
            :disabled="applyDisabled"
            class="button"
            type="info"
            size="small"
            @click="apply"
          >{{this.$t("invitation.apply")}}</van-button>
          <p
            v-show="requiredAmount"
            class="subtitle"
          >{{this.$t("invitation.reject_reason", {amount: this.requiredAmountApproximation })}}</p>
        </div>
      </van-tab>
      <van-tab :title="titleInvitees" name="invitees">
        <van-list v-model="loading" :finished="finished">
          <InvitationInviteeListItem
            class="layout"
            v-for="item in usedInvitations"
            :key="item.code"
            :invitation="item"
          />
        </van-list>
      </van-tab>
    </van-tabs>
  </div>
</template>

<script>
import InvitationCodeListItem from "../components/Invitation/InvitationCodeListItem";
import InvitationInviteeListItem from "../components/Invitation/InvitationInviteeListItem";
export default {
  name: "InvitationDetails",

  data() {
    return {
      titleCodes: this.$t("invitation.code"),
      titleInvitees: this.$t("invitation.invitees"),
      activeName: "codes",
      loading: false,
      finished: true,
      invitationsHistory: [],
      invitationsCurrent: [],
      ableToInvite: false,
      requiredAmount: null
    };
  },

  mounted() {
    this.GLOBAL.api.invitation
      .checkRule()
      .then(response => {
        if (response.data.pass) {
          this.ableToInvite = true;
        } else if (response.data.reject_reason && response.data.reject_reason.required_amount) {
          this.requiredAmount = response.data.reject_reason.required_amount;
        }
      });
    this.GLOBAL.api.invitation.index(false).then(response => {
      this.invitationsCurrent = response.data;
    });
    this.GLOBAL.api.invitation.index(true).then(response => {
      this.invitationsHistory = response.data;
    });
  },

  computed: {
    usedInvitations() {
      return this.invitationsHistory;
    },
    pendingInvitations() {
      return this.invitationsCurrent.filter(item => {
        return (
          item.is_used == true &&
          item.invitee &&
          item.invitee.state == "pending"
        );
      });
    },
    unusedInvitations() {
      return this.invitationsCurrent.filter(item => {
        return item.is_used == false;
      });
    },
    availableInvitations() {
      return this.unusedInvitations.concat(this.pendingInvitations) || [];
    },
    applyDisabled() {
      return !this.ableToInvite;
    },
    finishedText() {
      return this.invitationsCurrent.length === 0
        ? this.$t("invitation.no_code")
        : "";
    },
    requiredAmountApproximation() {
      return Number.parseFloat(this.requiredAmount).toFixed(8);
    }
  },

  components: {
    InvitationCodeListItem,
    InvitationInviteeListItem
  },

  methods: {
    apply() {
      this.GLOBAL.api.invitation
        .create()
        .then(response => {
          if (Array.isArray(response.data)) {
            this.invitationsCurrent.unshift(...response.data);
          } else if (!response.data.pass) {
            this.ableToInvite = false;
            if (response.data.reject_reason.requiredAmount != "") {
              this.requiredAmount = response.data.reject_reason.required_amount;
            }
          }
        });
    }
  }
};
</script>

<style lang="scss" scoped>
.layout {
  margin: 1rem;
  padding: 1rem;
}

.header {
  margin-top: 0;
  margin-bottom: 0;
}

.button {
  width: 100%;
  height: 3rem;
}

.subtitle {
  line-height: 22px;
  font-size: 14px;
  text-align: left;
  color: rgba(0, 0, 0, 0.25);
}

.pending {
  padding-left: 1rem;
}
</style>

