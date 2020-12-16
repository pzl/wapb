<template>
	<v-row>
		<v-col cols="12" md="">
			<v-alert v-if="alert" dense border="left" :type="alert.type" dismissable @input="alert = null">{{ alert.message }}</v-alert>
			<v-text-field name="url" label="URL" hint="URL Destination" v-model="toCreate.url" outlined />
			<create-meta
				ref="meta"
				v-bind.sync="toCreate"
			/>
			<v-btn block elevation="2" x-large color="success" :loading="isLoading" @click="create" :disabled="toCreate.url == ''">Create</v-btn>
			<instructions v-bind="instructions" />
		</v-col>
		<v-col cols="12" md="">
			<link-row v-for="l in links" :key="l.id" v-bind="l" />
		</v-col>
	</v-row>
</template>


<script>
import CreateMeta from '~/components/createMeta'
import LinkRow from '~/components/linkRow'
import Instructions from '~/components/instructions'

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
			links: [],

			instructions: {
				fields: [{ name: "url", type: "string", description: "(required) Destination URL" }],
				tools: [
					{
						name: "Curl",
						actions: [
							{
								name: "Create",
								methods: [
									{ name: "Simple", code: `curl -d "http://google.com" ${ this.$server }/api/v1/link` },
									{ name: "Simple with URL parameters", code: `curl -d "http://google.com" "${ this.$server }/api/v1/link?burn=true&ttl=3600"` },
									{ name: "Using Form Parameters", code: `curl -d "url=http://google.com" -d "burn=true" -d "ttl=3600" ${ this.$server }/api/v1/link` },
									{ name: "From your clipboard", code: `pbpaste | curl --data-binary @- ${ this.$server }/api/v1/link` },
									{ name: "JSON", code: `curl -d '{ "url": "http://google.com", "ttl":86500, "burn": false }' ${ this.$server }/api/v1/link` },
								]
							},
							{
								name: "Retrieve",
								note: `First, create an entry as above. When created successfully, you will get an ID back. It will either be at the end of a URL: (<code>${ this.$server }/link/fe4c71</code>) or as a JSON message: <code>{"id":"02b268", ... }</code>. We will use ID XXYY2 in our examples`,
								methods: [
									{ name: "As plain text", code: `curl -H "Accept: text/plain" ${ this.$server }/api/v1/link/XXYY2` },
									{ name: "To your clipboard!", code: `curl -H "Accept: text/plain" ${ this.$server }/api/v1/link/XXYY2 | pbcopy` },
									{ name: "As JSON", code: `curl ${ this.$server }/api/v1/link/XXYY2` },
								]
							},
							{ name: "Delete", methods: [{ code: `curl -X DELETE ${ this.$server }/api/v1/link/XXYY2` }] },
						]
					},
					{
						name: "HTTPie",
						actions: [
							{
								name: "Create",
								methods: [
									{ code: `http POST ${ this.$server }/api/v1/link 'url=http://google.com'` },
									{ code: `http POST ${ this.$server }/api/v1/link 'url=http://google.com' burn:=true ttl:=3600` },
									{ code: `echo '{ "url":"http://google.com", "burn":true }' | http POST ${ this.$server }/api/v1/link` },
								]
							},
							{ name: "Retrieve", methods: [{ code: `http ${ this.$server }/api/v1/link/XXYY2` }] },
							{ name: "Delete", methods: [{ code: `http DELETE ${ this.$server }/api/v1/link/XXYY2` }] },
						]
					}
				]
			}
		}
	},
	async asyncData({ $http, $server }) {
		const links = await $http.$get(`${$server}/api/v1/link`).then(d => d.data)
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
			const data = await this.$http.$post(`${this.$server}/api/v1/link`, this.toCreate).then(d => {
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
	components: { CreateMeta, LinkRow, Instructions }

}
</script>