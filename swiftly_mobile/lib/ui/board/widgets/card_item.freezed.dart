// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint
// ignore_for_file: unused_element, deprecated_member_use, deprecated_member_use_from_same_package, use_function_type_syntax_for_parameters, unnecessary_const, avoid_init_to_null, invalid_override_different_default_values_named, prefer_expression_function_bodies, annotate_overrides, invalid_annotation_target, unnecessary_question_mark

part of 'card_item.dart';

// **************************************************************************
// FreezedGenerator
// **************************************************************************

T _$identity<T>(T value) => value;

final _privateConstructorUsedError = UnsupportedError(
  'It seems like you constructed your class using `MyClass._()`. This constructor is only meant to be used by freezed and you are not supposed to need it nor use it.\nPlease check the documentation here for more information: https://github.com/rrousselGit/freezed#adding-getters-and-methods-to-our-models',
);

/// @nodoc
mixin _$CardItem {
  String get id => throw _privateConstructorUsedError;
  String get title => throw _privateConstructorUsedError;
  String get description => throw _privateConstructorUsedError;
  DateTime get createdAt => throw _privateConstructorUsedError;
  Priority get priority => throw _privateConstructorUsedError;
  String get category => throw _privateConstructorUsedError;
  String get columnId => throw _privateConstructorUsedError;

  /// Create a copy of CardItem
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  $CardItemCopyWith<CardItem> get copyWith =>
      throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $CardItemCopyWith<$Res> {
  factory $CardItemCopyWith(CardItem value, $Res Function(CardItem) then) =
      _$CardItemCopyWithImpl<$Res, CardItem>;
  @useResult
  $Res call({
    String id,
    String title,
    String description,
    DateTime createdAt,
    Priority priority,
    String category,
    String columnId,
  });
}

/// @nodoc
class _$CardItemCopyWithImpl<$Res, $Val extends CardItem>
    implements $CardItemCopyWith<$Res> {
  _$CardItemCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  /// Create a copy of CardItem
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? id = null,
    Object? title = null,
    Object? description = null,
    Object? createdAt = null,
    Object? priority = null,
    Object? category = null,
    Object? columnId = null,
  }) {
    return _then(
      _value.copyWith(
            id:
                null == id
                    ? _value.id
                    : id // ignore: cast_nullable_to_non_nullable
                        as String,
            title:
                null == title
                    ? _value.title
                    : title // ignore: cast_nullable_to_non_nullable
                        as String,
            description:
                null == description
                    ? _value.description
                    : description // ignore: cast_nullable_to_non_nullable
                        as String,
            createdAt:
                null == createdAt
                    ? _value.createdAt
                    : createdAt // ignore: cast_nullable_to_non_nullable
                        as DateTime,
            priority:
                null == priority
                    ? _value.priority
                    : priority // ignore: cast_nullable_to_non_nullable
                        as Priority,
            category:
                null == category
                    ? _value.category
                    : category // ignore: cast_nullable_to_non_nullable
                        as String,
            columnId:
                null == columnId
                    ? _value.columnId
                    : columnId // ignore: cast_nullable_to_non_nullable
                        as String,
          )
          as $Val,
    );
  }
}

/// @nodoc
abstract class _$$CardItemImplCopyWith<$Res>
    implements $CardItemCopyWith<$Res> {
  factory _$$CardItemImplCopyWith(
    _$CardItemImpl value,
    $Res Function(_$CardItemImpl) then,
  ) = __$$CardItemImplCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call({
    String id,
    String title,
    String description,
    DateTime createdAt,
    Priority priority,
    String category,
    String columnId,
  });
}

/// @nodoc
class __$$CardItemImplCopyWithImpl<$Res>
    extends _$CardItemCopyWithImpl<$Res, _$CardItemImpl>
    implements _$$CardItemImplCopyWith<$Res> {
  __$$CardItemImplCopyWithImpl(
    _$CardItemImpl _value,
    $Res Function(_$CardItemImpl) _then,
  ) : super(_value, _then);

  /// Create a copy of CardItem
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? id = null,
    Object? title = null,
    Object? description = null,
    Object? createdAt = null,
    Object? priority = null,
    Object? category = null,
    Object? columnId = null,
  }) {
    return _then(
      _$CardItemImpl(
        id:
            null == id
                ? _value.id
                : id // ignore: cast_nullable_to_non_nullable
                    as String,
        title:
            null == title
                ? _value.title
                : title // ignore: cast_nullable_to_non_nullable
                    as String,
        description:
            null == description
                ? _value.description
                : description // ignore: cast_nullable_to_non_nullable
                    as String,
        createdAt:
            null == createdAt
                ? _value.createdAt
                : createdAt // ignore: cast_nullable_to_non_nullable
                    as DateTime,
        priority:
            null == priority
                ? _value.priority
                : priority // ignore: cast_nullable_to_non_nullable
                    as Priority,
        category:
            null == category
                ? _value.category
                : category // ignore: cast_nullable_to_non_nullable
                    as String,
        columnId:
            null == columnId
                ? _value.columnId
                : columnId // ignore: cast_nullable_to_non_nullable
                    as String,
      ),
    );
  }
}

/// @nodoc

class _$CardItemImpl implements _CardItem {
  const _$CardItemImpl({
    required this.id,
    required this.title,
    required this.description,
    required this.createdAt,
    required this.priority,
    required this.category,
    required this.columnId,
  });

  @override
  final String id;
  @override
  final String title;
  @override
  final String description;
  @override
  final DateTime createdAt;
  @override
  final Priority priority;
  @override
  final String category;
  @override
  final String columnId;

  @override
  String toString() {
    return 'CardItem(id: $id, title: $title, description: $description, createdAt: $createdAt, priority: $priority, category: $category, columnId: $columnId)';
  }

  @override
  bool operator ==(Object other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$CardItemImpl &&
            (identical(other.id, id) || other.id == id) &&
            (identical(other.title, title) || other.title == title) &&
            (identical(other.description, description) ||
                other.description == description) &&
            (identical(other.createdAt, createdAt) ||
                other.createdAt == createdAt) &&
            (identical(other.priority, priority) ||
                other.priority == priority) &&
            (identical(other.category, category) ||
                other.category == category) &&
            (identical(other.columnId, columnId) ||
                other.columnId == columnId));
  }

  @override
  int get hashCode => Object.hash(
    runtimeType,
    id,
    title,
    description,
    createdAt,
    priority,
    category,
    columnId,
  );

  /// Create a copy of CardItem
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  @pragma('vm:prefer-inline')
  _$$CardItemImplCopyWith<_$CardItemImpl> get copyWith =>
      __$$CardItemImplCopyWithImpl<_$CardItemImpl>(this, _$identity);
}

abstract class _CardItem implements CardItem {
  const factory _CardItem({
    required final String id,
    required final String title,
    required final String description,
    required final DateTime createdAt,
    required final Priority priority,
    required final String category,
    required final String columnId,
  }) = _$CardItemImpl;

  @override
  String get id;
  @override
  String get title;
  @override
  String get description;
  @override
  DateTime get createdAt;
  @override
  Priority get priority;
  @override
  String get category;
  @override
  String get columnId;

  /// Create a copy of CardItem
  /// with the given fields replaced by the non-null parameter values.
  @override
  @JsonKey(includeFromJson: false, includeToJson: false)
  _$$CardItemImplCopyWith<_$CardItemImpl> get copyWith =>
      throw _privateConstructorUsedError;
}
