<template>
  <loading :loading="loading" :fullscreen="true">
    <div class="prepare-packet-page">
      <nav-bar
        :title="$t('prepare_packet.title')"
        :hasTopRight="false"
        :hasBack="true"
      ></nav-bar>
      <van-cell-group title>
        <row-select
          :index="0"
          :title="$t('prepare_packet.select_assets')"
          :columns="assets"
          placeholder="Tap to Select"
          @change="onChangeAsset"
        >
          <span slot="text">{{
            selectedAsset ? selectedAsset.text : "Tap to Select"
          }}</span>
        </row-select>
        <van-cell>
          <van-field
            type="number"
            v-model="form.amount"
            :label="$t('prepare_packet.amount')"
            :placeholder="
              $t('prepare_packet.placeholder_amount', { min: minAmount })
            "
          >
            <span slot="right-icon">{{
              selectedAsset ? selectedAsset.symbol : ""
            }}</span>
          </van-field>
        </van-cell>
        <van-cell title=" " :value="esitmatedValue"></van-cell>
        <van-cell>
          <van-field
            type="number"
            v-model="form.shares"
            :label="$t('prepare_packet.shares')"
            :placeholder="
              $t('prepare_packet.placeholder_shares', {
                count: Math.min(participantsCount, packetMaxCount)
              })
            "
          ></van-field>
        </van-cell>
        <van-cell>
          <van-field
            v-model="form.memo"
            rows="3"
            autosize
            type="textarea"
            :label="$t('prepare_packet.memo')"
            :placeholder="
              $t('prepare_packet.placeholder_memo', { max: maxMemoLength })
            "
            :maxlength="maxMemoLength"
          ></van-field>
        </van-cell>
      </van-cell-group>
      <van-row style="padding: 20px">
        <van-col span="24">
          <van-button
            style="width: 100%"
            type="info"
            :disabled="!validated"
            @click="pay"
            >{{ $t("prepare_packet.pay") }}</van-button
          >
        </van-col>
      </van-row>
    </div>
  </loading>
</template>

<script>
import NavBar from "@/components/Nav";
import RowSelect from "@/components/RowSelect";
import Row from "@/components/Nav";
import Loading from "@/components/Loading";
import uuid from "uuid";
import { Toast } from "vant";
import { CLIENT_ID } from "@/constants";

export default {
  name: "Prepare-Packet",
  props: {
    msg: String
  },
  data() {
    return {
      loading: false,
      coversationId: "",
      participantsCount: 0,
      assets: [],
      selectedAsset: null,
      packetMinAmountBase: "0.0001",
      packetMaxCount: 1000,
      form: {
        amount: "",
        shares: "",
        memo: this.$t("prepare_packet.default_memo", { symbol: "BTC" })
      }
    };
  },
  components: {
    NavBar,
    RowSelect,
    Loading
  },
  async mounted() {
    this.loading = true;
    let confInfo = await this.GLOBAL.api.website.config();
    if (confInfo && confInfo.data) {
      this.packetMinAmountBase = confInfo.data.redpacket_min_amount_base;
      this.packetMaxCount = confInfo.data.redpacket_max_count;
    }

    let prepareInfo = await this.GLOBAL.api.account.assets("redpacket");
    if (prepareInfo) {
      this.assets = prepareInfo.data.assets.map(x => {
        x.text = `${x.symbol} (${x.balance})`;
        return x;
      });
      if (this.assets.length) {
        this.selectedAsset = this.assets[0];
        this.form.memo = this.$t("prepare_packet.default_memo", {
          symbol: this.selectedAsset.symbol
        });
      }
      this.coversationId = prepareInfo.data.conversation.coversation_id;
      this.participantsCount = prepareInfo.data.conversation.participants_count;
    }
    this.loading = false;
  },
  computed: {
    maxMemoLength() {
      return 512;
    },
    validated() {
      if (
        this.form.amount &&
        this.form.shares &&
        this.selectedAsset &&
        parseFloat(this.form.amount) >= parseFloat(this.minAmount) &&
        this.form.shares <=
          Math.min(this.packetMaxCount, this.participantsCount)
      ) {
        return true;
      }
      return false;
    },
    esitmatedValue() {
      let val = 0;
      if (this.form.amount && this.selectedAsset) {
        val = this.form.amount * this.selectedAsset.price_usd;
      }
      return "≈$" + val.toLocaleString();
    },
    minAmount() {
      const base = parseFloat(this.packetMinAmountBase); // 1 usd
      if (base < 0.0001) {
        return 0;
      }
      if (this.selectedAsset) {
        return (base / this.selectedAsset.price_usd).toFixed(8);
      }
      return 0;
    }
  },
  methods: {
    async pay() {
      let payload = {
        amount: this.form.amount,
        total_count: parseInt(this.form.shares),
        greeting: this.form.memo,
        conversation_id: uuid.v4(),
        asset_id: this.selectedAsset.asset_id
      };

      if (payload.greeting.length > 512) {
        this.loading = false;
        Toast(this.$t("prepare_packet.memo_too_long_err"));
        return;
      }

      this.loading = true;
      let createResp = "";
      try {
        createResp = await this.GLOBAL.api.packet.create(payload);
        if (createResp.error) {
          this.loading = false;
          Toast("Error: " + createResp.error.description);
          return;
        }
      } catch (err) {
        this.loading = false;
        Toast("Error: " + err.description ? err.description : err.toString());
      }

      let pkt = createResp.data;
      setTimeout(() => {
        this.waitForPayment(pkt.packet_id);
      }, 1000);
      window.location.href = `mixin://pay?recipient=${CLIENT_ID}&asset=${
        this.selectedAsset.asset_id
      }&amount=${this.form.amount}&trace=${
        pkt.packet_id
      }&memo=${encodeURIComponent("SuperGroup LuckyCoin")}`;
    },
    onChangeAsset(ix) {
      this.selectedAsset = this.assets[ix];
      this.form.memo = this.$t("prepare_packet.default_memo", {
        symbol: this.selectedAsset.symbol
      });
    },
    async waitForPayment(packetId) {
      let resp = await this.GLOBAL.api.packet.show(packetId);
      if (resp.error) {
        setTimeout(() => {
          this.waitForPayment(packetId);
        }, 1500);
        return;
      }
      var data = resp.data;
      switch (data.state) {
        case "INITIAL":
          setTimeout(() => {
            this.waitForPayment(packetId);
          }, 1500);
          break;
        case "PAID":
        case "EXPIRED":
        case "REFUNDED":
          this.loading = false;
          this.$router.push("/packets/" + packetId);
          break;
      }
    }
  }
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.prepare-packet-page {
  padding-top: 60px;
}
h3 {
  margin: 40px 0 0;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
a {
  color: #42b983;
}
</style>
