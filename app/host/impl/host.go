package impl

import (
	"context"
	"github.com/infraboard/mcube/sqlbuilder"
	"github.com/puzhihao/gin-gorm-demo/app/host"
	"log"
)

func (i *impl) CreateHost(ctx context.Context, ins *host.Host) (*host.Host, error) {
	//校验数据合法性
	err := ins.Validate()
	if err != nil {
		return nil, err
	}

	//一次需要入库2个表，我们需要开启数据库事物
	//初始化一个事务
	tx, err := i.db.BeginTx(ctx, nil)
	if err != nil {
		print(err)
		return nil, err
	}

	//判断事务执行过程中是否有异常
	//有异常就回滚，无异常就提交
	//先执行prepare，防止sql注入
	resStmt, err := tx.Prepare(insertResourceSQL)
	if err != nil {
		return nil, err
	}
	_, err = resStmt.Exec(
		ins.Id, ins.Vendor, ins.Region, ins.Zone, ins.CreateAt, ins.ExpireAt, ins.Category, ins.Type, ins.InstanceId, ins.Name, ins.Description, ins.Status, ins.UpdateAt, ins.SyncAt, ins.SyncAccount, ins.PublicIP, ins.PrivateIP, ins.PayType, ins.DescribeHash, ins.ResourceHash,
	)
	if err != nil {
		return nil, err
	}
	defer resStmt.Close()

	hostStmt, err := tx.Prepare(insertHostSQL)
	if err != nil {
		return nil, err
	}
	defer hostStmt.Close()
	_, err = hostStmt.Exec(
		ins.Id, ins.CPU, ins.Memory, ins.GPUAmount, ins.GPUSpec, ins.OSType, ins.OSName, ins.SerialNumber, ins.ImageID, ins.InternetMaxBandwidthOut, ins.InternetMaxBandwidthIn, ins.KeyPairName, ins.SecurityGroups,
	)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			err := tx.Rollback()
			log.Printf("tx rollback error,%s", err)
		} else {
			err := tx.Commit()
			if err != nil {
				log.Printf("tx commit error,%s", err)
			}

		}
	}()
	return ins, nil
}

func (i *impl) QueryHost(ctx context.Context, req *host.QueryHostRequest) (*host.Set, error) {

	query := sqlbuilder.NewQuery(queryHostSQL).Order("create_at").Desc().Limit(int64((req.PageSize-1)*req.PageNumber), uint(req.PageNumber))
	//build查询语句
	querySQL, args := query.BuildQuery()
	log.Printf("sql:%s ,args: %v,", querySQL, args)
	stmt, err := i.db.Prepare(querySQL)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}

	//初始化要返回的数据
	set := host.NewSet()
	for rows.Next() {
		ins := host.NewDefaultHost()

		err := rows.Scan(
			&ins.Id, &ins.Vendor, &ins.Region, &ins.Zone, &ins.CreateAt, &ins.ExpireAt, &ins.Category, &ins.Type, &ins.InstanceId, &ins.Name, &ins.Description, &ins.Status, &ins.UpdateAt, &ins.SyncAt, &ins.SyncAccount, &ins.PublicIP, &ins.PrivateIP, &ins.PayType, &ins.DescribeHash, &ins.ResourceHash, &ins.Id, &ins.CPU, &ins.Memory, &ins.GPUAmount, &ins.GPUSpec, &ins.OSType, &ins.OSName, &ins.SerialNumber, &ins.ImageID, &ins.InternetMaxBandwidthOut, &ins.InternetMaxBandwidthIn, &ins.KeyPairName, &ins.SecurityGroups,
		)
		if err != nil {
			return nil, err
		}

		set.ADD(ins)
	}

	set.Total = int64(len(set.Items))
	////获取总数量
	////build count语句
	//countStr, countArgs := query.BuildCount()
	//fmt.Println(countArgs)
	//
	//log.Printf("sql:%s args: %v,", countStr, countArgs)
	//countStmt, err := i.db.Prepare(countStr)
	//if err != nil {
	//	return nil, err
	//}
	//defer countStmt.Close()
	//err = countStmt.QueryRow(countArgs...).Scan(&set.Total)
	//fmt.Println(set.Total)
	//if err != nil {
	//	return nil, err
	//}

	return set, nil
}
