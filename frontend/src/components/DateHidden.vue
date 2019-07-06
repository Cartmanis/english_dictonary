<template>
    <!-- :return-value.sync="dateFormatted" -->
    <v-menu ref = 'menu' :close-on-content-click="false" v-model="menuDate"
            :nudge-right="getNudge"
            lazy transition="scale-transition" offset-y full-width min-width="290px">
        <v-text-field slot="activator" prepend-icon="event" mask="##.##.####"
                      v-model="dateFormatted" :label= "labelDate" :readonly="readonly"
                      :rules ="[rules.required, rules.dateFormat, rules.dateValid]"
                      @input="onChange">
        </v-text-field>
        <v-date-picker v-model="date" locale="ru-ru" :color = "color">
            <v-spacer></v-spacer>
            <v-btn flat :color="color" @click="menuDate = false">Отмена</v-btn>
            <v-btn flat :color="color" @click="onClickOk">ОК</v-btn>
        </v-date-picker>
    </v-menu>
</template>
<script>
    export default {
        props: {
            label: {
                type: String
            },
            color: {
                type: String
            },
            readonly: {
                type: Boolean,
                default: false
            },
            required: {
                type: Boolean,
                default: false,
            }
        },
        data () {
            return {
                menuDate: false,
                date: null,
                dateFormatted: null,
                rules: {
                    required:  value =>  (!!value || !this.required || 'поле не должно быть пустым'),
                    dateFormat: value => {
                        if (value) {
                            return this.isFormatDate(value) || 'Формат даты должен быть: ДД.ММ.ГГГГ'
                        }
                        return true
                    },
                    dateValid: value => {
                        if (value) {
                            return this.isDateValid(value) || 'Не существующая дата'
                        }
                        return true
                    }
                }
            }
        },
        computed: {
            getRequired() {
                if (this.required) {
                    return '*'
                }
                return ''
            },
            labelDate() {
                if (this.label) {
                    return this.label + this.getRequired
                }
                return 'Some Date' + this.getRequired
            },
            getNudge () {
                return this.readonly ? "33" : "120"
            }
        },
        methods: {
            onClickOk () {
                this.onChange()
            },
            formatDate (date) {
                if (!date) return null
                const [year, month, day] = date.split('-')
                // .
                return `${day}${month}${year}`
            },
            onChange () {
                if (this.isFormatDate(this.dateFormatted) && this.isDateValid(this.dateFormatted)) {
                    // const [day, month, year] = this.dateFormatted.split('.')
                    const day = this.dateFormatted.slice(0, 2)
                    const month = this.dateFormatted.slice(2, 4)
                    const year = this.dateFormatted.slice(4)
                    this.date = `${year}-${month}-${day}`
                    this.$refs.menu.save(this.dateFormatted)
                    // this.dateFormatted
                    this.$emit('change-date', `${day}.${month}.${year}`)
                } else {
                    this.$emit('change-date', null)
                }
            },
            isDateValid (date) {
                if (!date) return false
                // const [day, month, year] = date.split('.')
                const day = date.slice(0, 2)
                const month = date.slice(2, 4)
                const year = date.slice(4)
                if (day && month && year) {
                    try {
                        const nDay = +day
                        const nMonth = +month - 1
                        const nYear = +year
                        const countDayCurrentMonth = 32 - new Date(nYear, nMonth, 32).getDate()
                        return countDayCurrentMonth >= nDay && (+month <= 12 && +month > 0) &&
                            (nYear <= 2100 && nYear >= 1800)
                    } catch (e) {
                        return false
                    }
                }
            },
            isFormatDate (date) {
                // const pattern = /([0-9]{2}).([0-9]{2}).([0-9]{4})/
                const patternMask = /[0-9]{8}/
                return patternMask.test(date) && date.length <= 10
            }
        },
        watch: {
            date (val) {
                this.dateFormatted = this.formatDate(this.date)
            }
        }
    }
</script>


