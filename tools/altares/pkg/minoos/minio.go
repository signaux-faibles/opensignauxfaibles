package minoos

import (
	"context"
	"io"
	"log/slog"
	"time"

	"github.com/jaswdr/faker"
	"github.com/minio/minio-go/v7"

	"opensignauxfaibles/tools/altares/pkg/utils"
)

const BUCKET_NAME = "bidon"
const REGION = "cloudgouv-eu-west-1"

var fake faker.Faker

type MC struct {
	delegate   *minio.Client
	bucketName string
	ctx        context.Context
}

func NewWithClient(minioClient *minio.Client) MC {
	bucketName := BUCKET_NAME
	mc := MC{minioClient, bucketName, context.Background()}
	mc.ensureBucketExists()
	return mc
}

func (m MC) ListAltaresFiles() []string {
	// liste les fichiers
	opts := minio.ListObjectsOptions{
		WithMetadata: true,
		Prefix:       "altares",
		Recursive:    true,
	}
	r := []string{}
	objects := m.delegate.ListObjects(m.ctx, m.bucketName, opts)
	for current := range objects {
		r = append(r, current.Key)
	}
	return r
}

func (m MC) GetAltaresFile(name string) *minio.Object {
	remoteName := "altares/" + name
	r, err := m.delegate.GetObject(m.ctx, m.bucketName, remoteName, minio.GetObjectOptions{Checksum: true})
	utils.ManageError(err, "erreur à la récupération du fichier "+name)
	slog.Debug("récupère le fichier", slog.String("name", remoteName), slog.Any("object", r))
	return r
}

func (m MC) PutAltaresFile(name string, reader io.Reader) (int64, time.Time) {
	slog.Info("pousse le fichier", slog.String("status", "start"), slog.String("filename", name))
	opts := minio.PutObjectOptions{
		UserMetadata:    map[string]string{"type": "fichier stock", "name": name},
		SendContentMd5:  true,
		ContentLanguage: "fr",
		ContentEncoding: "gzip",
		ContentType:     "binary/octet-stream",
	}
	info, err := m.delegate.PutObject(m.ctx, m.bucketName, "altares/"+name, reader, -1, opts)
	utils.ManageError(err, "erreur au poussage du fichier "+name)
	slog.Info(
		"pousse le fichier",
		slog.String("status", "end"),
		slog.Group("file",
			slog.String("key", info.Key),
			slog.Any("wrote", info.Size),
			slog.Any("lastModified", info.LastModified),
		),
	)
	return info.Size, info.LastModified
}

func (mc MC) ensureBucketExists() {
	exists, err := mc.delegate.BucketExists(mc.ctx, mc.bucketName)
	utils.ManageError(err, "erreur pendant la vérification de l'existence du bucket "+mc.bucketName)
	slog.Debug("le bucket existe", slog.String("name", mc.bucketName), slog.Bool("exists", exists))
	if !exists {
		slog.Info("crée le bucket", slog.String("name", mc.bucketName))
		options := minio.MakeBucketOptions{
			Region:        REGION,
			ObjectLocking: false,
		}
		err = mc.delegate.MakeBucket(mc.ctx, mc.bucketName, options)
		utils.ManageError(err, "erreur à la création du bucket "+mc.bucketName)
	}
	mc.configureBucket()
}

func (mc MC) configureBucket() {
	versioningEnabled := minio.BucketVersioningConfiguration{
		Status: "Enabled",
	}
	err := mc.delegate.SetBucketVersioning(mc.ctx, mc.bucketName, versioningEnabled)
	utils.ManageError(err, "erreur à l'activation du versioning")
	slog.Debug("versioning activé", slog.String("name", mc.bucketName))
}

//func (mc MC) CleanupVersionedBucket() {
//	slog.Warn("nettoyage du seau versionné", slog.String("status", "start"))
//	doneCh := make(chan struct{})
//	defer close(doneCh)
//	for obj := range mc.delegate.ListObjects(mc.ctx, mc.bucketName, minio.ListObjectsOptions{WithVersions: true, Recursive: true}) {
//		utils.ManageError(obj.Err, "erreur au listing des objets du seau")
//		if obj.Key != "" {
//			err := mc.delegate.RemoveObject(mc.ctx, mc.bucketName, obj.Key,
//				minio.RemoveObjectOptions{VersionID: obj.VersionID, GovernanceBypass: true})
//			utils.ManageError(err, "erreur à la suppression de l'objet "+obj.Key)
//		}
//	}
//	for objPartInfo := range mc.delegate.ListIncompleteUploads(context.Background(), mc.bucketName, "", true) {
//		utils.ManageError(objPartInfo.Err, "erreur au listing des poussages incomplets")
//		if objPartInfo.Key != "" {
//			err := mc.delegate.RemoveIncompleteUpload(context.Background(), mc.bucketName, objPartInfo.Key)
//			utils.ManageError(err, "erreur à la suppression du poussage incomplet de "+objPartInfo.Key)
//		}
//	}
//	// objects are already deleted, clear the buckets now
//	err := mc.delegate.RemoveBucket(context.Background(), mc.bucketName)
//	if err != nil {
//		for obj := range mc.delegate.ListObjects(context.Background(), mc.bucketName, minio.ListObjectsOptions{WithVersions: true, Recursive: true}) {
//			log.Println("found", obj.Key, obj.VersionID)
//		}
//		utils.ManageError(err, "erreur à la suppression du seau "+mc.bucketName)
//	}
//	slog.Warn("nettoyage du seau versionné", slog.String("status", "end"))
//}

//func generateIncrementCSV() string {
//	headers := []string{
//		"Référence Client",
//		"Siren",
//		"Siret",
//		"Raison sociale 1",
//		"Raison sociale 2",
//		"Enseigne",
//		"Sigle",
//		"Complément d'adresse",
//		"Adresse",
//		"Distribution spéciale",
//		"Code postal et bureau distributeur",
//		"Pays",
//		"Code postal",
//		"Ville",
//		"Qualité Etablissement",
//		"Code type d'établissement",
//		"Libellé type d'établissement",
//		"Etat d'activité établissement",
//		"Etat d'activité entreprise",
//		"Etat de procédure collective",
//		"Diffusible",
//		"Paydex",
//		"Retard moyen de paiements (j)",
//		"Nombre de fournisseurs analysés",
//		"Montant total des encours étudiés (€)",
//		"Montant total des encours échus non réglés (€)",
//		"FPI 30+",
//		"FPI 90+",
//		"Code du mouvement",
//		"Libellé du mouvement",
//		"Date d'effet du mouvement",
//	}
//	return ""
//}
