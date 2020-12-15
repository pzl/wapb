<template>
	<v-row>
		<v-col>
			<v-alert v-if="alert" dense border="left" :type="alert.type" dismissable @input="alert = null">{{ alert.message }}</v-alert>
			<v-text-field name="url" label="URL" hint="URL Destination" v-model="toCreate.url" outlined />
			<create-meta
				ref="meta"
				v-bind.sync="toCreate"
			/>
			<v-btn block elevation="2" x-large color="success" :loading="isLoading" @click="create" :disabled="toCreate.url == ''">Create</v-btn>
		</v-col>
		<v-col>
			<link-row v-for="l in links" :key="l.id" v-bind="l" />
		</v-col>
	</v-row>
</template>


<script>
import CreateMeta from '~/components/createMeta'
import LinkRow from '~/components/linkRow'

const createFactory = () => {
	return {
		url: '',
		burn: false,
		ttl: null,
		hidden: false,
	}
} 

export default {
	data () {
		return {
			toCreate: createFactory(),
			isLoading: false,
			alert: null,
			links: []
		}
	},
	async asyncData({ $http }) {
		const links = await $http.$get('http://localhost:7473/api/v1/link').then(d => d.data)
		return { links }
	},
	methods: {
		async create() {
			if (this.toCreate.url == '') {
				this.err = { message: 'Cannot create empty url' };
				return
			}
			this.alert = null;
			this.isLoading = true;
			const data = await this.$http.$post('http://localhost:7473/api/v1/link', this.toCreate).then(d => {
				this.isLoading = false;
				this.toCreate = createFactory()
				this.$refs.meta.hideTTL()
				this.links.push(d)
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
		},
	},
	components: { CreateMeta, LinkRow }

}
</script>