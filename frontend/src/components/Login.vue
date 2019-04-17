<template>
  <v-container fluid fill-height>
    <v-layout align-center wrap justify-center>
        <v-alert
          style="position: absolute; top: 0px; z-index: 10;"
          v-model="loginError"
          dismissible
          type="error"
        >
          Identification incorrecte
        </v-alert>
      <v-flex class="login" xs12 sm6 md3>
        <span><v-img class="center" max-width="100px" src="/static/logo_signaux_faibles_cercle.svg"></v-img></span>
        <h1><span class="fblue">Signaux</span>·<span class="fred">Faibles</span></h1>
      </v-flex>
      <v-flex class="login" xs12 sm6 md3>
          <v-form @submit="login">
          <v-text-field  prepend-icon="person" label="Adresse électronique" v-model="email"></v-text-field >
          <v-text-field  prepend-icon="lock" type="password" :label="passwordHint" v-model="password"></v-text-field>
          <v-btn round color="primary" type="submit">Connexion</v-btn><br/>
          <div v-if="browserToken">
            <v-dialog
              persistent
              height="400px"
              width="500"
              v-model="passwordHelp"
            >
              <a slot="activator" style="font-size: 10px" href='#'>Mot de passe oublié ?</a>
                <v-card >
                  <v-toolbar color="#eef" class="toolbar-widget elevation-1" dense>
                    <v-icon>mdi-lock-question</v-icon>
                    <v-spacer></v-spacer>
                    <span style="font-size: 15px; font-weight: 600">
                      Récupération du compte
                    </span>
                  </v-toolbar>
                  <div v-if="state==1">
                    <v-card-text>
                      <b>Etape 1</b> · Etape 2 · Fin<br/>
                      Saisissez votre adresse électronique, un code de récupération vous sera envoyé.
                    </v-card-text>
                    <v-card-text>
                      <v-text-field
                      prepend-icon="mdi-at"
                      label="Adresse E-Mail"
                      :rules="[rules.required, rules.email]"
                      v-model="recoveryEmail">
                      </v-text-field >
                    </v-card-text>
                    <v-card-actions>
                      <v-spacer></v-spacer>
                      <v-btn
                        color="primary"
                        flat
                        @click="passwordHelp = false"
                      >
                        annuler
                      </v-btn>
                      <v-btn
                        color="primary"
                        flat
                        :disabled="!validEmail"
                        @click="getRecovery()"
                      >
                        suivant
                      </v-btn>
                    </v-card-actions>
                  </div>
                  <div v-if="state==2">
                    <v-card-text>
                      Etape 1 · <b>Etape 2</b> · Fin<br/>
                      Saisissez le code de récupération reçu dans votre boite aux lettres électronique
                    </v-card-text>
                    <v-card-text>
                      <v-text-field length=6 solo prepend-icon="mdi-pound" label="Code de récupération" v-model="recoveryCode"></v-text-field >
                    </v-card-text>
                      <v-card-actions>
                      <v-spacer></v-spacer>
                      <v-btn
                        color="primary"
                        flat
                        @click="rollbackRecovery()"
                      >
                        recommencer
                      </v-btn>
                      <v-btn
                        color="primary"
                        flat
                        @click="checkRecovery()"
                      >
                        suivant
                      </v-btn>
                    </v-card-actions>
                  </div>
                  <div v-if="state==3">
                    <v-card-text>
                      Etape 1 · Etape 2 · <b>Fin</b><br/>
                      Saisissez votre nouveau mot de passe
                    </v-card-text>
                    <v-card-text>
                      <v-text-field
                        prepend-icon="fa-key"
                        type="password"
                        label="Nouveau mot de passe"
                        :rules="[rules.complexity]"
                        v-model="newPassword">
                      </v-text-field >
                      <v-progress-linear :color="passwordStrengthColor" v-model="passwordStrength"></v-progress-linear>
                      Mot de passe {{passwordStrengthComment}}
                    </v-card-text>
                    <v-card-actions>
                      <v-spacer></v-spacer>
                      <v-btn
                        color="primary"
                        flat
                        @click="rollbackRecovery()"
                      >
                        recommencer
                      </v-btn>
                      <v-btn
                      color="primary"
                      flat
                      :disabled="passwordStrength<77"
                      @click="setPasswordRecovery()"
                      >
                        terminer
                      </v-btn>
                    </v-card-actions>
                  </div>
                </v-card>
              </v-dialog>
            </div>
          </v-form>
        <v-dialog
        v-model="verifDialog"
        width="500"
        >
          <v-card>
            <v-toolbar color="#eef" class="toolbar-widget elevation-1" dense>
              <v-icon>mdi-lock</v-icon>
              <v-spacer></v-spacer>
              Validation de votre connexion
            </v-toolbar>
            C'est la première fois que vous vous identifiez sur ce navigateur, merci de confirmer cette connexion en saisissant le code que vous allez recevoir à l'adresse {{ email }}.
            <v-text-field length=6 solo prepend-icon="mdi-pound" label="Code de vérification" v-model="loginCode"></v-text-field >
            <v-card-actions>
              <v-spacer></v-spacer>
              <v-btn
              @click="checkLogin"
              >se connecter</v-btn>
            </v-card-actions>
          </v-card>
        </v-dialog>
      </v-flex>
    </v-layout>
  </v-container>
</template>

<script>
import zxcvbn from 'zxcvbn'
export default {
  data () {
    return {
      passwordHelp: false,
      recoveryEmail: null,
      recoveryCode: null,
      newPassword: null,
      newPasswordConfirm: null,
      state: 1,
      rules: {
        required: value => !!value || 'Obligatoire.',
        email: value => this.validEmail || 'Adresse invalide.',
        samePassword: value => this.samePassword || 'Mot de passe différents'
      },
      verifDialog: false,
      loginCode: ''

    }
  },
  computed: {
    passwordHint () {
      switch (this.loginTry) {
        case 4:
          this.passwordHelp = true
          return 'Mot de passe'
        case 3:
          return 'Mot de passe'
        case 2:
          return 'Mot de passe: 2 essais restants'
        case 1:
          return 'Mot de passe: dernier essai'
      }
    },
    loginTry () {
      return this.$store.state.loginTry
    },
    loginError () {
      return this.$store.state.loginError
    },
    validEmail () {
      const pattern = /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/
      return pattern.test(this.recoveryEmail)
    },
    samePassword () {
      return this.newPassword === this.newPasswordConfirm
    },
    passwordStrength () {
      return zxcvbn(this.newPassword || '').score * 25 + 2
    },
    passwordStrengthComment () {
      if (this.passwordStrength < 50) {
        return 'insuffisant'
      } else if (this.passwordStrength < 75) {
        return 'faible'
      } else if (this.passwordStrength < 100) {
        return 'acceptable'
      } else {
        return 'résistant'
      }
    },
    passwordStrengthColor () {
      if (this.passwordStrength < 50) {
        return 'red'
      } else if (this.passwordStrength < 75) {
        return 'orange'
      } else if (this.passwordStrength < 100) {
        return 'yellow'
      } else {
        return 'green'
      }
    },
    email: {
      get () {
        return this.$store.state.credentials.email
      },
      set (email) {
        this.$store.commit('setEmail', email)
      }
    },
    password: {
      get (password) {
        return this.$store.state.credentials.password
      },
      set (password) {
        this.$store.commit('setPassword', password)
      }
    },
    token () {
      return this.$store.state.token
    },
    browserToken () {
      return this.$localStore.state.browserToken
    }
  },
  methods: {
    checkLogin () {
      this.verifDialog = false
      this.$store.dispatch('checkLogin', this.loginCode)
    },
    login () {
      // if (this.browserToken != null) {
      this.$store.dispatch('login')
      // } else {
      //   this.$store.dispatch('getLogin')
      //   this.verifDialog = true
      // }
    },
    getRecovery () {
      let self = this
      let parameters = {
        email: this.recoveryEmail,
        browserToken: this.browserToken
      }
      this.$axios.post('/login/recovery/get', parameters).then(r => {
        self.state = 2
      })
    },
    checkRecovery () {
      this.state = 3
    },
    setPasswordRecovery () {
      let parameters = {
        email: this.recoveryEmail,
        browserToken: this.browserToken,
        code: this.recoveryCode,
        password: this.newPassword
      }
      this.$axios.post('/login/recovery/setPassword', parameters).then(r => {
        this.passwordHelp = false
        this.email = this.recoveryEmail
        this.password = this.newPassword
        this.$store.dispatch('login')
      })
    },
    rollbackRecovery () {
      this.state = 1
    }
  }
}
</script>

<style scoped>
  .login {
        text-align: center;     /* will center text in <p>, which is not a flex item */
        align-items: center
  }
  h1 {
    font-family: 'Quicksand', sans-serif;
  }
  .center {
    display: block;
    margin: 0 auto;
  }
  #toolbar {
    font-size: 25px;
    color: white;
  }
  span.fblue {
    font-family: 'Quicksand', sans-serif;
    color: #20459a
  }
  span.fred {
    font-family: 'Quicksand', sans-serif;
    color: #e9222e
  }
  .box {
    background-color: #fff;
  }
</style>
