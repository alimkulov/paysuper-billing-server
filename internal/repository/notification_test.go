package repository

import (
	"context"
	"errors"
	"github.com/bxcodec/faker"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/paysuper/paysuper-billing-server/internal/config"
	"github.com/paysuper/paysuper-billing-server/internal/mocks"
	"github.com/paysuper/paysuper-billing-server/internal/repository/models"
	"github.com/paysuper/paysuper-proto/go/billingpb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	mongodb "gopkg.in/paysuper/paysuper-database-mongo.v2"
	"testing"
)

type NotificationTestSuite struct {
	suite.Suite
	db         mongodb.SourceInterface
	repository *notificationRepository
	log        *zap.Logger
}

func Test_NotificationTestSuite(t *testing.T) {
	suite.Run(t, new(NotificationTestSuite))
}

func (suite *NotificationTestSuite) SetupTest() {
	models.InitFakeCustomProviders()

	_, err := config.NewConfig()
	assert.NoError(suite.T(), err, "Config load failed")

	suite.log, err = zap.NewProduction()
	assert.NoError(suite.T(), err, "Logger initialization failed")

	suite.db, err = mongodb.NewDatabase()
	assert.NoError(suite.T(), err, "Database connection failed")

	suite.repository = &notificationRepository{db: suite.db, cache: &mocks.CacheInterface{}, mapper: models.NewNotificationMapper()}
}

func (suite *NotificationTestSuite) TearDownTest() {

	if err := suite.db.Drop(); err != nil {
		suite.FailNow("Database deletion failed", "%v", err)
	}

	if err := suite.db.Close(); err != nil {
		suite.FailNow("Database close failed", "%v", err)
	}
}


func (suite *NotificationTestSuite) Test_NewNotificationRepository_Ok() {
	repository := NewNotificationRepository(suite.db)
	assert.IsType(suite.T(), &notificationRepository{}, repository)
}

func (suite *NotificationTestSuite) Test_NotificationRepositoryInsert_Ok() {
	shouldBe := require.New(suite.T())

	obj := &billingpb.Notification{}
	shouldBe.NoError(faker.FakeData(obj))
	err := suite.repository.Insert(context.TODO(), obj)
	shouldBe.NoError(err)

	obj2, err := suite.repository.GetById(context.TODO(), obj.Id)
	shouldBe.NoError(err)

	assert.ObjectsAreEqualValues(obj2, obj)
}

func (suite *NotificationTestSuite) Test_NotificationRepositoryInsert_MapperError() {
	shouldBe := require.New(suite.T())

	mapper := &mocks.Mapper{}
	mapper.On("MapMgoToObject", mock.Anything).Return(nil, errors.New("error"))
	mapper.On("MapObjectToMgo", mock.Anything).Return(nil, errors.New("error"))

	suite.repository.mapper = mapper

	obj := &billingpb.Notification{}
	shouldBe.NoError(faker.FakeData(obj))
	err := suite.repository.Insert(context.TODO(), obj)
	shouldBe.Error(err)
}

func (suite *NotificationTestSuite) Test_NotificationRepositoryInsert_BadData() {
	shouldBe := require.New(suite.T())

	obj := &billingpb.Notification{}
	shouldBe.NoError(faker.FakeData(obj))
	obj.MerchantId = "asdasd"
	err := suite.repository.Insert(context.TODO(), obj)
	shouldBe.Error(err)

	obj = &billingpb.Notification{}
	shouldBe.NoError(faker.FakeData(obj))
	obj.CreatedAt = &timestamp.Timestamp{Seconds: -100000000000000}
	err = suite.repository.Insert(context.TODO(), obj)
	shouldBe.Error(err)
}

func (suite *NotificationTestSuite) Test_NotificationRepositoryFind_MapperError() {
	shouldBe := require.New(suite.T())

	obj := &billingpb.Notification{}
	shouldBe.NoError(faker.FakeData(obj))
	obj.IsSystem = false
	err := suite.repository.Insert(context.TODO(), obj)
	shouldBe.NoError(err)

	mapper := &mocks.Mapper{}
	mapper.On("MapMgoToObject", mock.Anything).Return(nil, errors.New("error"))
	mapper.On("MapObjectToMgo", mock.Anything).Return(nil, errors.New("error"))

	suite.repository.mapper = mapper

	_, err = suite.repository.Find(context.TODO(), obj.MerchantId, obj.UserId, 0, nil, 0, 100)
	shouldBe.Error(err)
}

func (suite *NotificationTestSuite) Test_NotificationRepositoryGetById_MapperError() {
	shouldBe := require.New(suite.T())

	mapper := &mocks.Mapper{}
	mapper.On("MapMgoToObject", mock.Anything).Return(nil, errors.New("error"))
	mapper.On("MapObjectToMgo", mock.Anything).Return(nil, errors.New("error"))

	suite.repository.mapper = mapper

	_, err := suite.repository.GetById(context.TODO(), primitive.NewObjectID().Hex())
	shouldBe.Error(err)
}


func (suite *NotificationTestSuite) Test_NotificationRepositoryGetById_NotFoundError() {
	shouldBe := require.New(suite.T())

	_, err := suite.repository.GetById(context.TODO(), primitive.NewObjectID().Hex())
	shouldBe.Error(err)
}


func (suite *NotificationTestSuite) Test_NotificationRepositoryGetById_IdError() {
	shouldBe := require.New(suite.T())

	_, err := suite.repository.GetById(context.TODO(), "HelloThere")
	shouldBe.Error(err)
}

func (suite *NotificationTestSuite) Test_NotificationRepositoryFind_Ok() {
	shouldBe := require.New(suite.T())

	obj := &billingpb.Notification{}
	shouldBe.NoError(faker.FakeData(obj))
	obj.IsSystem = false
	err := suite.repository.Insert(context.TODO(), obj)
	shouldBe.NoError(err)

	obj2 := &billingpb.Notification{}
	shouldBe.NoError(faker.FakeData(obj2))
	obj2.UserId = obj.UserId
	obj2.MerchantId = obj.MerchantId
	obj2.IsSystem = true
	err = suite.repository.Insert(context.TODO(), obj2)
	shouldBe.NoError(err)

	notifs, err := suite.repository.Find(context.TODO(), obj.MerchantId, obj.UserId, 0, nil, 0, 100)
	shouldBe.NoError(err)
	shouldBe.Len(notifs, 2)

	notifs, err = suite.repository.Find(context.TODO(), obj2.MerchantId, obj2.UserId, 1, nil, 0, 100)
	shouldBe.NoError(err)
	shouldBe.Len(notifs, 1)

	notifs, err = suite.repository.Find(context.TODO(), obj2.MerchantId, obj2.UserId, 2, nil, 0, 100)
	shouldBe.NoError(err)
	shouldBe.Len(notifs, 1)
}

func (suite *NotificationTestSuite) Test_NotificationRepositoryFindCount_Ok() {
	shouldBe := require.New(suite.T())

	obj := &billingpb.Notification{}
	shouldBe.NoError(faker.FakeData(obj))
	obj.IsSystem = false
	err := suite.repository.Insert(context.TODO(), obj)
	shouldBe.NoError(err)

	notifs, err := suite.repository.FindCount(context.TODO(), obj.MerchantId, obj.UserId, 0)
	shouldBe.NoError(err)
	shouldBe.EqualValues(notifs, 1)
}