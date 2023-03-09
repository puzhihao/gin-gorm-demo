package impl

import (
	"context"
	"errors"
	"github.com/infraboard/mcube/sqlbuilder"
	"github.com/puzhihao/gin-gorm-demo/app/host"
	"github.com/rs/xid"
	"log"
	"time"
)

func (i *impl) CreateHost(ctx context.Context, ins *host.Host) (*host.Host, error) {
	//校验数据合法性
	err := ins.Validate()
	if err != nil {
		return nil, err
	}

	// 分布式ID, app, instance, ip, mac, ......, idc(), region
	ins.Resource.Id = xid.New().String()
	if ins.Resource.CreateAt == 0 {
		ins.Resource.CreateAt = time.Now().UnixMilli()
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

	if req.Keywords != "" {
		query.Where(" r.name LIKE ?", "%"+req.Keywords+"%")
	}
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

func (i *impl) DescribeHost(ctx context.Context, req *host.DescribeHostRequest) (*host.Host, error) {

	query := sqlbuilder.NewQuery(queryHostSQL).Where("r.id=?", req.Id)

	sqlStr, args := query.BuildQuery()
	log.Printf("sql:%s ,args: %v,", sqlStr, args)
	stmt, err := i.db.Prepare(sqlStr)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	ins := host.NewDefaultHost()
	err = stmt.QueryRow(args...).Scan(
		&ins.Id, &ins.Vendor, &ins.Region, &ins.Zone, &ins.CreateAt, &ins.ExpireAt, &ins.Category, &ins.Type, &ins.InstanceId, &ins.Name, &ins.Description, &ins.Status, &ins.UpdateAt, &ins.SyncAt, &ins.SyncAccount, &ins.PublicIP, &ins.PrivateIP, &ins.PayType, &ins.DescribeHash, &ins.ResourceHash, &ins.Id, &ins.CPU, &ins.Memory, &ins.GPUAmount, &ins.GPUSpec, &ins.OSType, &ins.OSName, &ins.SerialNumber, &ins.ImageID, &ins.InternetMaxBandwidthOut, &ins.InternetMaxBandwidthIn, &ins.KeyPairName, &ins.SecurityGroups,
	)
	if err != nil {
		return nil, err
	}

	return ins, nil

}

func (i *impl) UpdateHost(ctx context.Context, req *host.UpdateHostRequest) (*host.Host, error) {
	//查询出需要更新的资源信息
	hostDesc := host.DescribeHostRequest{
		Id: req.Id,
	}
	ins, err := i.DescribeHost(ctx, &hostDesc)
	if err != nil {
		return nil, err
	}

	//对象更新（PATCH｜PUT）
	switch req.UpdateMode {
	case host.PUT:
		//全量更新
		ins.Update(req.Resource, req.DescribeHost)

	case host.PATCH:
		//部分更新
		err := ins.Patch(req.Resource, req.DescribeHost)
		if err != nil {
			return nil, err
		}
	}

	stmtRES, err := i.db.Prepare(updateResourceSQL)
	if err != nil {
		return nil, err
	}
	defer stmtRES.Close()
	_, err = stmtRES.Exec(ins.Vendor, ins.Region, ins.Zone, ins.Description, ins.Id)
	if err != nil {
		return nil, err
	}
	//校验更新后是否合法
	err = ins.Validate()
	if err != nil {
		return nil, err
	}

	stmtHOST, err := i.db.Prepare(updateHostSQL)
	if err != nil {
		return nil, err
	}
	defer stmtHOST.Close()
	_, err = stmtHOST.Exec(ins.CPU, ins.Memory, ins.Id)
	if err != nil {
		return nil, err
	}

	ins, err = i.DescribeHost(ctx, &hostDesc)
	if err != nil {
		return nil, err
	}
	//fmt.Println(res)
	//fmt.Println(hostRes)
	return ins, err

}

func (i *impl) DeleteHost(ctx context.Context, req *host.DeleteHostRequest) (*host.Set, error) {

	dres, err := i.db.Prepare(deleteResourceSQL)

	if err != nil {
		return nil, err
	}
	defer dres.Close()
	resStr, err := dres.Exec(req.Id)
	if err != nil {
		return nil, err
	}

	dhost, err := i.db.Prepare(deleteHostSQL)
	if err != nil {
		return nil, err
	}
	defer dhost.Close()
	hostStr, err := dhost.Exec(req.Id)
	if err != nil {
		return nil, err
	}
	query := host.QueryHostRequest{
		PageNumber: 20,
		PageSize:   1,
	}
	queryHost, err := i.QueryHost(ctx, &query)

	affected, _ := resStr.RowsAffected()
	rowsAffected, _ := hostStr.RowsAffected()

	if affected+rowsAffected == 0 {
		return nil, errors.New("主机不存在")
	}

	return queryHost, nil
}
