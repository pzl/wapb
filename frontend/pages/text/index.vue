<template>
	<v-row>
		<v-col>
			<v-alert v-if="alert" dense border="left" :type="alert.type" dismissable @input="alert = null">{{ alert.message }}</v-alert>
			<v-textarea name="text" label="Text" hint="New Text" v-model="toCreate.text" outlined />
			<v-row no-gutter>
				<v-col>
					<v-checkbox name="Expire" label="Expire After Time" v-model="showTTL" />
				</v-col>
				<v-col>
					<v-select v-if="showTTL" v-model="selectedTimePreset" :items="timePresets" :hint="TTLHint" persistent-hint return-object single-line />
					<v-text-field v-if="showTTL && selectedTimePreset.value == -1" v-model="toCreate.ttl" hide-details type="number" label="Expire After (seconds)" />
				</v-col>
			</v-row>
			<v-checkbox v-model="toCreate.burn" label="Burn After Reading" />
			<v-checkbox v-model="toCreate.hidden" label="Hidden" />
			<v-btn block elevation="2" x-large color="success" :loading="isLoading" @click="create" :disabled="toCreate.text == ''">Create</v-btn>
		</v-col>
		<v-col>
			<div v-for="t in texts" :key="t.id">
				<p><nuxt-link :to="'/text/'+t.id">{{ t.text }}</nuxt-link></p>
			</div>
		</v-col>
	</v-row>
</template>


<script>

const createFactory = () => {
	return {
		text: '',
		burn: false,
		ttl: null,
		hidden: false,
	}
} 

export default {
	data () {
		return {
			toCreate: createFactory(),
			showTTL: false,
			timePresets: [
				{ text: '1 hour', value: 60*60 },
				{ text: '5 hours', value: 5*60*60 },
				{ text: '1 day', value: 24*60*60 },
				{ text: '1 week', value: 7*24*60*60 },
				{ text: 'Custom', value: -1 },
			],
			selectedTimePreset: {},
			isLoading: false,
			alert: null,
			texts: []
		}
	},
	async asyncData() {
		const texts = await fetch('http://localhost:7473/api/v1/text').then(res => res.json()).then(d => d.data)
		return { texts }
	},
	methods: {
		async create() {
			if (this.toCreate.text == '') {
				this.err = { message: 'Cannot create empty text' };
				return
			}
			this.alert = null;
			this.isLoading = true;
			const data = await this.$http.$post('http://localhost:7473/api/v1/text', this.toCreate).then(d => {
				this.isLoading = false;
				this.toCreate = createFactory()
				this.showTTL = false;
				this.texts.push(d)
				this.alert = {
					type: "success",
					message: `Created ${d.id}`,
				}
			}).catch(e => {
				this.isLoading = false;
				this.alert = {
					type: "error",
					message: e.message,
				}
			})
		}
	},
	computed: {
		TTLHint() {
			if (this.selectedTimePreset && this.selectedTimePreset.value >= 0) {
				return `${this.selectedTimePreset.value}s`;
			}
			return ''
		}
	},
	watch: {
		selectedTimePreset(newVal, oldVal) {
			if (newVal.value < 0) {
				this.toCreate.ttl = 0;
			} else {
				this.toCreate.ttl = newVal.value;				
			}
		},
		showTTL(newVal, oldVal) {
			if (!newVal) {
				this.toCreate.ttl = null;
			} else if (this.selectedTimePreset.value) {
				this.toCreate.ttl = this.selectedTimePreset.value;
			}
		}
	}
}
</script>