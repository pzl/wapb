<template>
	<v-container>
		<v-row no-gutters>
			<v-col>
				<v-checkbox name="Expire" label="Expire After Time" v-model="showTTL" />
			</v-col>
			<v-col>
				<v-select persistent-hint return-object single-line
					v-if="showTTL"
					v-model="selectedTimePreset"
					:items="timePresets"
					:hint="TTLHint" 
				/>
				<v-text-field type="number" label="Expire After (seconds)" hide-details 
					v-if="showTTL && selectedTimePreset.value == -1"
					@change="$emit('update:ttl', parseInt($event))"
				/>
			</v-col>
		</v-row>
		<v-checkbox @change="$emit('update:burn', $event)" label="Burn After Reading" />
		<v-checkbox @change="$emit('update:hidden', $event)" label="Hidden" />
	</v-container>
</template>


<script>
export default {
	props: {
		hidden: {}, // prevent the "hidden" property from hiding the actual HTML
	},
	data() {
		return {
			showTTL: false,
			timePresets: [
				{ text: '1 hour', value: 60*60 },
				{ text: '5 hours', value: 5*60*60 },
				{ text: '1 day', value: 24*60*60 },
				{ text: '1 week', value: 7*24*60*60 },
				{ text: 'Custom', value: -1 },
			],
			selectedTimePreset: {},
		}
	},
	methods: {
		hideTTL() {
			this.showTTL = false;
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
				this.$emit('update:ttl', 0)
			} else {
				this.$emit('update:ttl', newVal.value)				
			}
		},
		showTTL(newVal, oldVal) {
			if (!newVal) {
				this.$emit('update:ttl', null);
			} else if (this.selectedTimePreset.value) {
				this.$emit('update:ttl', this.selectedTimePreset.value);
			}
		}
	}
}
</script>
